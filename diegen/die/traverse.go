/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package die

import (
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-tools/pkg/loader"
)

// codeWriter assists in writing out Go code lines and blocks to a writer.
type codeWriter struct {
	out io.Writer
}

// Linef writes a single line with formatting (as per fmt.Sprintf).
func (c *codeWriter) Linef(line string, args ...interface{}) {
	fmt.Fprintf(c.out, line+"\n", args...)
}

func (c *codeWriter) Doc(lines string) string {
	b := strings.Builder{}
	for i, line := range strings.Split(lines, "\n") {
		if i != 0 {
			b.WriteString("\n//\n")
		}
		b.WriteString("// " + strings.TrimSpace(line))
	}
	return b.String()
}

// importsList keeps track of required imports, automatically assigning aliases
// to import statement.
type importsList struct {
	byPath  map[string]string
	byAlias map[string]string

	pkg *loader.Package
}

func (l *importsList) AliasedRef(importPath, typeName string) string {
	if importPath == "" || l.pkg.ID == importPath {
		return typeName
	}
	alias := l.NeedImport(importPath)
	return fmt.Sprintf("%s.%s", alias, typeName)
}

// NeedImport marks that the given package is needed in the list of imports,
// returning the ident (import alias) that should be used to reference the package.
func (l *importsList) NeedImport(importPath string) string {
	// we get an actual path from Package, which might include venddored
	// packages if running on a package in vendor.
	if ind := strings.LastIndex(importPath, "/vendor/"); ind != -1 {
		importPath = importPath[ind+8: /* len("/vendor/") */]
	}

	// check to see if we've already assigned an alias, and just return that.
	alias, exists := l.byPath[importPath]
	if exists {
		return alias
	}

	// otherwise, calculate an import alias by joining path parts till we get something unique
	restPath, nextWord := path.Split(importPath)

	for otherPath, exists := "", true; exists && otherPath != importPath; otherPath, exists = l.byAlias[alias] {
		if restPath == "" {
			// do something else to disambiguate if we're run out of parts and
			// still have duplicates, somehow
			alias += "x"
		}

		// can't have a first digit, per Go identifier rules, so just skip them
		for firstRune, runeLen := utf8.DecodeRuneInString(nextWord); unicode.IsDigit(firstRune); firstRune, runeLen = utf8.DecodeRuneInString(nextWord) {
			nextWord = nextWord[runeLen:]
		}

		// make a valid identifier by replacing "bad" characters with underscores
		nextWord = strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				return r
			}
			return '_'
		}, nextWord)

		alias = nextWord + alias
		if len(restPath) > 0 {
			restPath, nextWord = path.Split(restPath[:len(restPath)-1] /* chop off final slash */)
		}
	}

	l.byPath[importPath] = alias
	l.byAlias[alias] = importPath
	return alias
}

// ImportSpecs returns a string form of each import spec
// (i.e. `alias "path/to/import").  Aliases are only present
// when they don't match the package name.
func (l *importsList) ImportSpecs() []string {
	res := make([]string, 0, len(l.byPath))
	for importPath, alias := range l.byPath {
		pkg := l.pkg.Imports()[importPath]
		if pkg != nil && pkg.Name == alias {
			// don't print if alias is the same as package name
			// (we've already taken care of duplicates).
			res = append(res, fmt.Sprintf("%q", importPath))
		} else {
			res = append(res, fmt.Sprintf("%s %q", alias, importPath))
		}
	}
	return res
}

// copyMethodMakers makes DeepCopy (and related) methods for Go types,
// writing them to its codeWriter.
type copyMethodMaker struct {
	pkg *loader.Package
	*importsList
	*codeWriter
	test *copyMethodMaker
	dies sets.Set[string]
}

func (c *copyMethodMaker) GenerateMethodsFor(die Die, fields []Field) {
	c.generateDieFor(die)
	c.generateObjectMethodsFor(die)
}

func (c *copyMethodMaker) GenerateFieldFor(field Field, die Die) {
	c.Linef("")
	if field.Doc != "" {
		c.Linef("%s", c.Doc(field.Doc))
	}
	c.Linef("func (d *%s) %s(v %s%s) *%s {", field.Receiver, field.Name, field.TypePrefix, c.AliasedRef(field.TypePackage, field.Type), die.Type)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		r.%s = v", field.Name)
	c.Linef("	})")
	c.Linef("}")

	if field.Type == "IntOrString" && field.TypePackage == "k8s.io/apimachinery/pkg/util/intstr" && (field.TypePrefix == "*" || field.TypePrefix == "") {
		c.Linef("")
		c.Linef("func (d *%s) %sInt(i int) *%s {", field.Receiver, field.Name, die.Type)
		c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
		c.Linef("		v := %s(i)", c.AliasedRef(field.TypePackage, "FromInt"))
		if field.TypePrefix == "*" {
			c.Linef("		r.%s = &v", field.Name)
		} else {
			c.Linef("		r.%s = v", field.Name)
		}
		c.Linef("	})")
		c.Linef("}")

		c.Linef("")
		c.Linef("func (d *%s) %sString(s string) *%s {", field.Receiver, field.Name, die.Type)
		c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
		c.Linef("		v := %s(s)", c.AliasedRef(field.TypePackage, "FromString"))
		if field.TypePrefix == "*" {
			c.Linef("		r.%s = &v", field.Name)
		} else {
			c.Linef("		r.%s = v", field.Name)
		}
		c.Linef("	})")
		c.Linef("}")
	}
}

func (c *copyMethodMaker) generateDieFor(die Die) {
	c.Linef("")
	c.Linef("var %s = (&%s{}).DieFeed(%s{})", die.Blank, die.Type, c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("")
	if die.Doc != "" {
		c.Linef("%s", c.Doc(die.Doc))
	}
	c.Linef("type %s struct {", die.Type)
	if die.Object {
		c.Linef("	%s", c.AliasedRef("reconciler.io/dies/apis/meta/v1", "FrozenObjectMeta"))
	}
	c.Linef("	mutable bool")
	c.Linef("	r %s", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("}")
	c.Linef("")

	c.generateDieImmutableMethodFor(die)
	c.generateDieFeedMethodFor(die)
	c.generateDieReleaseMethodFor(die)
	c.generateDieStampMethodFor(die)
	c.generateDieWithMethodFor(die)
	c.generateDeepCopyMethodFor(die)

	c.test.generateMissingFieldTestFor(die)
}

func (c *copyMethodMaker) generateDieImmutableMethodFor(die Die) {
	c.Linef("")
	c.Linef("// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).")
	c.Linef("func (d *%s) DieImmutable(immutable bool) *%s {", die.Type, die.Type)
	c.Linef("	if d.mutable == !immutable {")
	c.Linef("		return d")
	c.Linef("	}")
	c.Linef("	d = d.DeepCopy()")
	c.Linef("	d.mutable = !immutable")
	c.Linef("	return d")
	c.Linef("}")
}

func (c *copyMethodMaker) generateDieFeedMethodFor(die Die) {
	c.Linef("")
	c.Linef("// DieFeed returns a new die with the provided resource.")
	c.Linef("func (d *%s) DieFeed(r %s) *%s {", die.Type, c.AliasedRef(die.TargetPackage, die.TargetType), die.Type)
	c.Linef("	if d.mutable {")
	if die.Object {
		c.Linef("		d.FrozenObjectMeta = %s(r.ObjectMeta)", c.AliasedRef("reconciler.io/dies/apis/meta/v1", "FreezeObjectMeta"))
	}
	c.Linef("		d.r = r")
	c.Linef("		return d")
	c.Linef("	}")
	c.Linef("	return &%s{", die.Type)
	if die.Object {
		c.Linef("		FrozenObjectMeta: %s(r.ObjectMeta),", c.AliasedRef("reconciler.io/dies/apis/meta/v1", "FreezeObjectMeta"))
	}
	c.Linef("		mutable: d.mutable,")
	c.Linef("		r: r,")
	c.Linef("	}")
	c.Linef("}")

	c.Linef("")
	c.Linef("// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.")
	c.Linef("func (d *%s) DieFeedPtr(r *%s) *%s {", die.Type, c.AliasedRef(die.TargetPackage, die.TargetType), die.Type)
	c.Linef("	if r == nil {")
	c.Linef("		r = &%s{}", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	}")
	c.Linef("	return d.DieFeed(*r)")
	c.Linef("}")

	c.Linef("")
	c.Linef("// DieFeedJSON returns a new die with the provided JSON. Panics on error.")
	c.Linef("func (d *%s) DieFeedJSON(j []byte) *%s {", die.Type, die.Type)
	c.Linef("	r := %s{}", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	if err := %s(j, &r); err != nil {", c.AliasedRef("encoding/json", "Unmarshal"))
	c.Linef("		panic(err)")
	c.Linef("	}")
	c.Linef("	return d.DieFeed(r)")
	c.Linef("}")

	c.Linef("")
	c.Linef("// DieFeedYAML returns a new die with the provided YAML. Panics on error.")
	c.Linef("func (d *%s) DieFeedYAML(y []byte) *%s {", die.Type, die.Type)
	c.Linef("	r := %s{}", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	if err := %s(y, &r); err != nil {", c.AliasedRef("sigs.k8s.io/yaml", "Unmarshal"))
	c.Linef("		panic(err)")
	c.Linef("	}")
	c.Linef("	return d.DieFeed(r)")
	c.Linef("}")

	c.Linef("")
	c.Linef("// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.")
	c.Linef("func (d *%s) DieFeedYAMLFile(name string) *%s {", die.Type, die.Type)
	c.Linef("	y, err := %s(name)", c.AliasedRef("os", "ReadFile"))
	c.Linef("	if err != nil {")
	c.Linef("		panic(err)")
	c.Linef("	}")
	c.Linef("	return d.DieFeedYAML(y)")
	c.Linef("}")

	c.Linef("")
	c.Linef("// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.")
	c.Linef("func (d *%s) DieFeedRawExtension(raw %s) *%s {", die.Type, c.AliasedRef("k8s.io/apimachinery/pkg/runtime", "RawExtension"), die.Type)
	c.Linef("	j, err := %s(raw)", c.AliasedRef("encoding/json", "Marshal"))
	c.Linef("	if err != nil {")
	c.Linef("		panic(err)")
	c.Linef("	}")
	c.Linef("	return d.DieFeedJSON(j)")
	c.Linef("}")
}

func (c *copyMethodMaker) generateDieReleaseMethodFor(die Die) {
	c.Linef("")
	c.Linef("// DieRelease returns the resource managed by the die.")
	c.Linef("func (d *%s) DieRelease() %s {", die.Type, c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	if d.mutable {")
	c.Linef("		return d.r")
	c.Linef("	}")
	c.Linef("	return *d.r.DeepCopy()")
	c.Linef("}")

	c.Linef("")
	c.Linef("// DieReleasePtr returns a pointer to the resource managed by the die.")
	c.Linef("func (d *%s) DieReleasePtr() *%s {", die.Type, c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	r := d.DieRelease()")
	c.Linef("	return &r")
	c.Linef("}")
	if die.Object {
		c.Linef("")
		c.Linef("// DieReleaseUnstructured returns the resource managed by the die as an unstructured object. Panics on error.")
		c.Linef("func (d *%s) DieReleaseUnstructured() *%s {", die.Type, c.AliasedRef("k8s.io/apimachinery/pkg/apis/meta/v1/unstructured", "Unstructured"))
		c.Linef("	r := d.DieReleasePtr()")
		c.Linef("	u, err := %s.ToUnstructured(r)", c.AliasedRef("k8s.io/apimachinery/pkg/runtime", "DefaultUnstructuredConverter"))
		c.Linef("	if err != nil {")
		c.Linef("		panic(err)")
		c.Linef("	}")
		c.Linef("	return &%s{", c.AliasedRef("k8s.io/apimachinery/pkg/apis/meta/v1/unstructured", "Unstructured"))
		c.Linef("		Object: u,")
		c.Linef("	}")
		c.Linef("}")
	}

	c.Linef("")
	c.Linef("// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.")
	c.Linef("func (d *%s) DieReleaseJSON() []byte {", die.Type)
	c.Linef("	r := d.DieReleasePtr()")
	c.Linef("	j, err := %s(r)", c.AliasedRef("encoding/json", "Marshal"))
	c.Linef("	if err != nil {")
	c.Linef("		panic(err)")
	c.Linef("	}")
	c.Linef("	return j")
	c.Linef("}")

	c.Linef("")
	c.Linef("// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.")
	c.Linef("func (d *%s) DieReleaseYAML() []byte {", die.Type)
	c.Linef("	r := d.DieReleasePtr()")
	c.Linef("	y, err := %s(r)", c.AliasedRef("sigs.k8s.io/yaml", "Marshal"))
	c.Linef("	if err != nil {")
	c.Linef("		panic(err)")
	c.Linef("	}")
	c.Linef("	return y")
	c.Linef("}")

	c.Linef("")
	c.Linef("// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.")
	c.Linef("func (d *%s) DieReleaseRawExtension() %s {", die.Type, c.AliasedRef("k8s.io/apimachinery/pkg/runtime", "RawExtension"))
	c.Linef("	j := d.DieReleaseJSON()")
	c.Linef("	raw := %s{}", c.AliasedRef("k8s.io/apimachinery/pkg/runtime", "RawExtension"))
	c.Linef("	if err := %s(j, &raw); err != nil {", c.AliasedRef("encoding/json", "Unmarshal"))
	c.Linef("		panic(err)")
	c.Linef("	}")
	c.Linef("	return raw")
	c.Linef("}")
}

func (c *copyMethodMaker) generateDieStampMethodFor(die Die) {
	c.Linef("")
	c.Linef("// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.")
	c.Linef("func (d *%s) DieStamp(fn func(r *%s)) *%s {", die.Type, c.AliasedRef(die.TargetPackage, die.TargetType), die.Type)
	c.Linef("	r := d.DieRelease()")
	c.Linef("	fn(&r)")
	c.Linef("	return d.DieFeed(r)")
	c.Linef("}")

	c.Linef("")
	c.Linef("// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.")
	c.Linef("//")
	c.Linef("// Future iterations will improve type coercion from the resource to the callback argument.")
	c.Linef("func (d *%s) DieStampAt(jp string, fn interface{}) *%s {", die.Type, die.Type)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		if ni := %s(fn).Type().NumIn(); ni != 1 {", c.AliasedRef("reflect", "ValueOf"))
	c.Linef("			panic(%s(\"callback function must have 1 input parameters, found %%d\", ni))", c.AliasedRef("fmt", "Errorf"))
	c.Linef("		}")
	c.Linef("		if no := %s(fn).Type().NumOut(); no != 0 {", c.AliasedRef("reflect", "ValueOf"))
	c.Linef("			panic(%s(\"callback function must have 0 output parameters, found %%d\", no))", c.AliasedRef("fmt", "Errorf"))
	c.Linef("		}")
	c.Linef("")
	c.Linef("		cp := %s(\"\")", c.AliasedRef("k8s.io/client-go/util/jsonpath", "New"))
	c.Linef("		if err := cp.Parse(%s(\"{%%s}\", jp)); err != nil {", c.AliasedRef("fmt", "Sprintf"))
	c.Linef("			panic(err)")
	c.Linef("		}")
	c.Linef("		cr, err := cp.FindResults(r)")
	c.Linef("		if err != nil {")
	c.Linef("			// errors are expected if a path is not found")
	c.Linef("			return")
	c.Linef("		}")
	c.Linef("		for _, cv := range cr[0] {")
	c.Linef("			arg0t := %s(fn).Type().In(0)", c.AliasedRef("reflect", "ValueOf"))
	c.Linef("")
	c.Linef("			var args []%s", c.AliasedRef("reflect", "Value"))
	c.Linef("			if cv.Type().AssignableTo(arg0t) {")
	c.Linef("				args = []%s{cv}", c.AliasedRef("reflect", "Value"))
	c.Linef("			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {")
	c.Linef("				args = []%s{cv.Addr()}", c.AliasedRef("reflect", "Value"))
	c.Linef("			} else {")
	c.Linef("				panic(%s(\"callback function must accept value of type %%q, found type %%q\", cv.Type(), arg0t))", c.AliasedRef("fmt", "Errorf"))
	c.Linef("			}")
	c.Linef("")
	c.Linef("			%s(fn).Call(args)", c.AliasedRef("reflect", "ValueOf"))
	c.Linef("		}")
	c.Linef("	})")
	c.Linef("}")
}

func (c *copyMethodMaker) generateDieWithMethodFor(die Die) {
	c.Linef("")
	c.Linef("// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.")
	c.Linef("func (d *%s) DieWith(fns ...func(d *%s)) *%s {", die.Type, die.Type, die.Type)
	c.Linef("	nd := %s.DieFeed(d.DieRelease()).DieImmutable(false)", die.Blank)
	c.Linef("	for _, fn := range fns {")
	c.Linef("		if fn != nil {")
	c.Linef("			fn(nd)")
	c.Linef("		}")
	c.Linef("	}")
	c.Linef("	return d.DieFeed(nd.DieRelease())")
	c.Linef("}")
}

func (c *copyMethodMaker) generateDeepCopyMethodFor(die Die) {
	c.Linef("")
	c.Linef("// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.")
	c.Linef("func (d *%s) DeepCopy() *%s {", die.Type, die.Type)
	c.Linef("	r := *d.r.DeepCopy()")
	c.Linef("	return &%s{", die.Type)
	if die.Object {
		c.Linef("		FrozenObjectMeta: %s(r.ObjectMeta),", c.AliasedRef("reconciler.io/dies/apis/meta/v1", "FreezeObjectMeta"))
	}
	c.Linef("		mutable: d.mutable,")
	c.Linef("		r: r,")
	c.Linef("	}")
	c.Linef("}")
}

func (c *copyMethodMaker) generateObjectMethodsFor(die Die) {
	if die.Object {
		c.generateRuntimeObjectMethodsFor(die)
		c.generateJSONMethodsFor(die)
		c.generateMetadataDieMethodFor(die)
		if c.dies.Has(die.Spec) {
			c.generateSpecMethodFor(die)
		}
		if c.dies.Has(die.Status) {
			c.generateStatusMethodFor(die)
		}
	}
}

func (c *copyMethodMaker) generateRuntimeObjectMethodsFor(die Die) {
	c.Linef("")
	c.Linef("var _ %s = (*%s)(nil)", c.AliasedRef("k8s.io/apimachinery/pkg/runtime", "Object"), die.Type)

	c.Linef("")
	c.Linef("func (d *%s) DeepCopyObject() %s {", die.Type, c.AliasedRef("k8s.io/apimachinery/pkg/runtime", "Object"))
	c.Linef("	return d.r.DeepCopy()")
	c.Linef("}")

	c.Linef("")
	c.Linef("func (d *%s) GetObjectKind() %s {", die.Type, c.AliasedRef("k8s.io/apimachinery/pkg/runtime/schema", "ObjectKind"))
	c.Linef("	r := d.DieRelease()")
	c.Linef("	return r.GetObjectKind()")
	c.Linef("}")
}

func (c *copyMethodMaker) generateJSONMethodsFor(die Die) {
	c.Linef("")
	c.Linef("func (d *%s) MarshalJSON() ([]byte, error) {", die.Type)
	c.Linef("	return %s(d.r)", c.AliasedRef("encoding/json", "Marshal"))
	c.Linef("}")

	c.Linef("")
	c.Linef("func (d *%s) UnmarshalJSON(b []byte) error {", die.Type)
	c.Linef("	if !d.mutable {")
	c.Linef("		return %s(\"cannot unmarshal into immutable dies, create a mutable version first\")", c.AliasedRef("fmt", "Errorf"))
	c.Linef("	}")
	c.Linef("	r := &%s{}", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	err := %s(b, r)", c.AliasedRef("encoding/json", "Unmarshal"))
	c.Linef("	*d = *d.DieFeed(*r)")
	c.Linef("	return err")
	c.Linef("}")
}

func (c *copyMethodMaker) generateMetadataDieMethodFor(die Die) {
	if die.APIVersion != "" || die.Kind != "" {
		c.Linef("")
		c.Linef("// DieDefaultTypeMetadata sets the APIVersion and Kind to %q and %q respectively.", die.APIVersion, die.Kind)
		c.Linef("func (d *%s) DieDefaultTypeMetadata() *%s {", die.Type, die.Type)
		c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
		c.Linef("		r.APIVersion = %q", die.APIVersion)
		c.Linef("		r.Kind = %q", die.Kind)
		c.Linef("	})")
		c.Linef("}")
	}

	c.Linef("")
	c.Linef("// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources")
	c.Linef("func (d *%s) APIVersion(v string) *%s {", die.Type, die.Type)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		r.APIVersion = v")
	c.Linef("	})")
	c.Linef("}")

	c.Linef("")
	c.Linef("// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds")
	c.Linef("func (d *%s) Kind(v string) *%s {", die.Type, die.Type)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		r.Kind = v")
	c.Linef("	})")
	c.Linef("}")

	c.Linef("")
	c.Linef("// TypeMetadata standard object's type metadata.")
	c.Linef("func (d *%s) TypeMetadata(v %s) *%s {", die.Type, c.AliasedRef("k8s.io/apimachinery/pkg/apis/meta/v1", "TypeMeta"), die.Type)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		r.TypeMeta = v")
	c.Linef("	})")
	c.Linef("}")

	c.Linef("")
	c.Linef("// TypeMetadataDie stamps the resource's TypeMeta field with a mutable die.")
	c.Linef("func (d *%s) TypeMetadataDie(fn func(d *%s)) *%s {", die.Type, c.AliasedRef("reconciler.io/dies/apis/meta/v1", "TypeMetaDie"), die.Type)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		d := %s.DieImmutable(false).DieFeed(r.TypeMeta)", c.AliasedRef("reconciler.io/dies/apis/meta/v1", "TypeMetaBlank"))
	c.Linef("		fn(d)")
	c.Linef("		r.TypeMeta = d.DieRelease()")
	c.Linef("	})")
	c.Linef("}")

	c.Linef("")
	c.Linef("// Metadata standard object's metadata.")
	c.Linef("func (d *%s) Metadata(v %s) *%s {", die.Type, c.AliasedRef("k8s.io/apimachinery/pkg/apis/meta/v1", "ObjectMeta"), die.Type)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		r.ObjectMeta = v")
	c.Linef("	})")
	c.Linef("}")

	c.Linef("")
	c.Linef("// MetadataDie stamps the resource's ObjectMeta field with a mutable die.")
	c.Linef("func (d *%s) MetadataDie(fn func(d *%s)) *%s {", die.Type, c.AliasedRef("reconciler.io/dies/apis/meta/v1", "ObjectMetaDie"), die.Type)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		d := %s.DieImmutable(false).DieFeed(r.ObjectMeta)", c.AliasedRef("reconciler.io/dies/apis/meta/v1", "ObjectMetaBlank"))
	c.Linef("		fn(d)")
	c.Linef("		r.ObjectMeta = d.DieRelease()")
	c.Linef("	})")
	c.Linef("}")
}

func (c *copyMethodMaker) generateSpecMethodFor(die Die) {
	c.Linef("")
	c.Linef("// SpecDie stamps the resource's spec field with a mutable die.")
	c.Linef("func (d *%s) SpecDie(fn func(d *%s)) *%s {", die.Type, die.SpecType, die.Type)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		d := %s.DieImmutable(false).DieFeed(r.Spec)", die.SpecBlank)
	c.Linef("		fn(d)")
	c.Linef("		r.Spec = d.DieRelease()")
	c.Linef("	})")
	c.Linef("}")
}

func (c *copyMethodMaker) generateStatusMethodFor(die Die) {
	c.Linef("")
	c.Linef("// StatusDie stamps the resource's status field with a mutable die.")
	c.Linef("func (d *%s) StatusDie(fn func(d *%s)) *%s {", die.Type, die.StatusType, die.Type)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		d := %s.DieImmutable(false).DieFeed(r.Status)", die.StatusBlank)
	c.Linef("		fn(d)")
	c.Linef("		r.Status = d.DieRelease()")
	c.Linef("	})")
	c.Linef("}")
}

func (c *copyMethodMaker) generateMissingFieldTestFor(die Die) {
	ignore := []string{}
	for _, f := range die.IgnoreFields {
		ignore = append(ignore, strconv.Quote(f))
	}

	c.Linef("")
	c.Linef("func Test%s_MissingMethods(t *%s) {", die.Type, c.AliasedRef("testing", "T"))
	c.Linef("	die := %s", c.AliasedRef(c.pkg.PkgPath, die.Blank))
	c.Linef("	ignore := []string{%s}", strings.Join(ignore, ", "))
	c.Linef("	diff := %s(die).Delete(ignore...)", c.AliasedRef("reconciler.io/dies/testing", "DieFieldDiff"))
	c.Linef("	if diff.Len() != 0 {")
	c.Linef("		t.Errorf(\"found missing fields for %s: %%s\", diff.List())", die.Type)
	c.Linef("	}")
	c.Linef("}")
}
