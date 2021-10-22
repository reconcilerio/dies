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
	dies sets.String
}

func (c *copyMethodMaker) GenerateMethodsFor(die Die, fields []Field) {
	c.generateInterface(die, fields)
	c.generateDieFor(die)
	c.generateObjectMethodsFor(die)
}

func (c *copyMethodMaker) GenerateFieldFor(field Field, die Die) {
	c.Linef("")
	c.Linef("func (d *%s) %s(v %s%s) %s {", field.Receiver, field.Name, field.TypePrefix, c.AliasedRef(field.TypePackage, field.Type), die.Interface)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		r.%s = v", field.Name)
	c.Linef("	})")
	c.Linef("}")

	if field.Type == "IntOrString" && field.TypePackage == "k8s.io/apimachinery/pkg/util/intstr" && (field.TypePrefix == "*" || field.TypePrefix == "") {
		c.Linef("")
		c.Linef("func (d *%s) %sInt(i int) %s {", field.Receiver, field.Name, die.Interface)
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
		c.Linef("func (d *%s) %sString(s string) %s {", field.Receiver, field.Name, die.Interface)
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

func (c *copyMethodMaker) generateInterface(die Die, fields []Field) {
	c.Linef("")
	c.Linef("type %s interface {", die.Interface)
	c.Linef("	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.")
	c.Linef("	DieStamp(fn func(r *%s)) %s", c.AliasedRef(die.TargetPackage, die.TargetType), die.Interface)
	c.Linef("	// DieFeed returns a new die with the provided resource.")
	c.Linef("	DieFeed(r %s) %s", c.AliasedRef(die.TargetPackage, die.TargetType), die.Interface)
	c.Linef("	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.")
	c.Linef("	DieFeedPtr(r *%s) %s", c.AliasedRef(die.TargetPackage, die.TargetType), die.Interface)
	c.Linef("	// DieRelease returns the resource managed by the die.")
	c.Linef("	DieRelease() %s", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	// DieReleasePtr returns a pointer to the resource managed by the die.")
	c.Linef("	DieReleasePtr() *%s", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).")
	c.Linef("	DieImmutable(immutable bool) %s", die.Interface)
	c.Linef("	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.")
	c.Linef("	DeepCopy() %s", die.Interface)
	c.Linef("")
	if die.CustomMethodInterface != "" {
		c.Linef("	%s", die.CustomMethodInterface)
	}
	if die.Object {
		c.Linef("	// MetadataDie stamps the resource's ObjectMeta field with a mutable die.")
		c.Linef("	MetadataDie(fn func(d %s)) %s", c.AliasedRef("dies.dev/apis/meta/v1", "ObjectMetaDie"), die.Interface)
		if c.dies.Has(die.SpecName) {
			c.Linef("	// SpecDie stamps the resource's spec field with a mutable die.")
			c.Linef("	SpecDie(fn func(d %s)) %s", die.SpecInterface, die.Interface)
		}
		if c.dies.Has(die.StatusName) {
			c.Linef("	// StatusDie stamps the resource's status field with a mutable die.")
			c.Linef("	StatusDie(fn func(d %s)) %s", die.StatusInterface, die.Interface)
		}
	}
	for _, f := range fields {
		if f.Doc != "" {
			c.Linef("%s", c.Doc(f.Doc))
		}
		c.Linef("	%s(v %s%s) %s", f.Name, f.TypePrefix, c.AliasedRef(f.TypePackage, f.Type), die.Interface)
	}
	if die.Object {
		c.Linef("")
		c.Linef("	%s", c.AliasedRef("k8s.io/apimachinery/pkg/runtime", "Object"))
		c.Linef("	%s", c.AliasedRef("k8s.io/apimachinery/pkg/apis/meta/v1", "Object"))
		c.Linef("	%s", c.AliasedRef("k8s.io/apimachinery/pkg/apis/meta/v1", "ObjectMetaAccessor"))
	}

	c.Linef("}")
}

func (c *copyMethodMaker) generateDieFor(die Die) {
	c.Linef("")
	c.Linef("var _ %s = (*%s)(nil)", die.Interface, die.Type)
	c.Linef("var %s = (&%s{}).DieFeed(%s{})", die.Blank, die.Type, c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("")
	if die.Doc != "" {
		c.Linef("%s", c.Doc(die.Doc))
	}
	c.Linef("type %s struct {", die.Type)
	if die.Object {
		c.Linef("	%s", c.AliasedRef("dies.dev/apis/meta/v1", "FrozenObjectMeta"))
	}
	c.Linef("	mutable bool")
	c.Linef("	r %s", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("}")
	c.Linef("")

	c.generateDieImmutableMethodFor(die)
	c.generateDieFeedMethodFor(die)
	c.generateDieReleaseMethodFor(die)
	c.generateDieStampMethodFor(die)
	c.generateDeepCopyMethodFor(die)

	c.test.generateMissingFieldTestFor(die)
}

func (c *copyMethodMaker) generateDieImmutableMethodFor(die Die) {
	c.Linef("")
	c.Linef("func (d *%s) DieImmutable(immutable bool) %s {", die.Type, die.Interface)
	c.Linef("	if d.mutable == !immutable {")
	c.Linef("		return d")
	c.Linef("	}")
	c.Linef("	d = d.DeepCopy().(*%s)", die.Type)
	c.Linef("	d.mutable = !immutable")
	c.Linef("	return d")
	c.Linef("}")
}

func (c *copyMethodMaker) generateDieFeedMethodFor(die Die) {
	c.Linef("")
	c.Linef("func (d *%s) DieFeed(r %s) %s {", die.Type, c.AliasedRef(die.TargetPackage, die.TargetType), die.Interface)
	c.Linef("	if d.mutable {")
	if die.Object {
		c.Linef("		d.FrozenObjectMeta = %s(r.ObjectMeta)", c.AliasedRef("dies.dev/apis/meta/v1", "FreezeObjectMeta"))
	}
	c.Linef("		d.r = r")
	c.Linef("		return d")
	c.Linef("	}")
	c.Linef("	return &%s{", die.Type)
	if die.Object {
		c.Linef("		FrozenObjectMeta: %s(r.ObjectMeta),", c.AliasedRef("dies.dev/apis/meta/v1", "FreezeObjectMeta"))
	}
	c.Linef("		mutable: d.mutable,")
	c.Linef("		r: r,")
	c.Linef("	}")
	c.Linef("}")

	c.Linef("")
	c.Linef("func (d *%s) DieFeedPtr(r *%s) %s {", die.Type, c.AliasedRef(die.TargetPackage, die.TargetType), die.Interface)
	c.Linef("	if r == nil {")
	c.Linef("		r = &%s{}", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	}")
	c.Linef("	return d.DieFeed(*r)")
	c.Linef("}")
}

func (c *copyMethodMaker) generateDieReleaseMethodFor(die Die) {
	c.Linef("")
	c.Linef("func (d *%s) DieRelease() %s {", die.Type, c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	if d.mutable {")
	c.Linef("		return d.r")
	c.Linef("	}")
	c.Linef("	return *d.r.DeepCopy()")
	c.Linef("}")
	c.Linef("")
	c.Linef("func (d *%s) DieReleasePtr() *%s {", die.Type, c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	r := d.DieRelease()")
	c.Linef("	return &r")
	c.Linef("}")
}

func (c *copyMethodMaker) generateDieStampMethodFor(die Die) {
	c.Linef("")
	c.Linef("func (d *%s) DieStamp(fn func(r *%s)) %s {", die.Type, c.AliasedRef(die.TargetPackage, die.TargetType), die.Interface)
	c.Linef("	r := d.DieRelease()")
	c.Linef("	fn(&r)")
	c.Linef("	return d.DieFeed(r)")
	c.Linef("}")
}

func (c *copyMethodMaker) generateDeepCopyMethodFor(die Die) {
	c.Linef("")
	c.Linef("func (d *%s) DeepCopy() %s {", die.Type, die.Interface)
	c.Linef("	r := *d.r.DeepCopy()")
	c.Linef("	return &%s{", die.Type)
	if die.Object {
		c.Linef("		FrozenObjectMeta: %s(r.ObjectMeta),", c.AliasedRef("dies.dev/apis/meta/v1", "FreezeObjectMeta"))
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
		if c.dies.Has(die.SpecName) {
			c.generateSpecMethodFor(die)
		}
		if c.dies.Has(die.StatusName) {
			c.generateStatusMethodFor(die)
		}
	}
}

func (c *copyMethodMaker) generateRuntimeObjectMethodsFor(die Die) {
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
	c.Linef("	if d == %s {", die.Blank)
	c.Linef("		return %s(\"cannot unmarshal into the root object, create a copy first\")", c.AliasedRef("fmt", "Errorf"))
	c.Linef("	}")
	c.Linef("	r := &%s{}", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("	err := %s(b, r)", c.AliasedRef("encoding/json", "Unmarshal"))
	c.Linef("	*d = *d.DieFeed(*r).(*%s)", die.Type)
	c.Linef("	return err")
	c.Linef("}")
}

func (c *copyMethodMaker) generateMetadataDieMethodFor(die Die) {
	c.Linef("")
	c.Linef("func (d *%s) MetadataDie(fn func(d %s)) %s {", die.Type, c.AliasedRef("dies.dev/apis/meta/v1", "ObjectMetaDie"), die.Interface)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		d := %s.DieImmutable(false).DieFeed(r.ObjectMeta)", c.AliasedRef("dies.dev/apis/meta/v1", "ObjectMetaBlank"))
	c.Linef("		fn(d)")
	c.Linef("		r.ObjectMeta = d.DieRelease()")
	c.Linef("	})")
	c.Linef("}")
}

func (c *copyMethodMaker) generateSpecMethodFor(die Die) {
	c.Linef("")
	c.Linef("func (d *%s) SpecDie(fn func(d %s)) %s {", die.Type, die.SpecInterface, die.Interface)
	c.Linef("	return d.DieStamp(func(r *%s) {", c.AliasedRef(die.TargetPackage, die.TargetType))
	c.Linef("		d := %s.DieImmutable(false).DieFeed(r.Spec)", die.SpecBlank)
	c.Linef("		fn(d)")
	c.Linef("		r.Spec = d.DieRelease()")
	c.Linef("	})")
	c.Linef("}")
}

func (c *copyMethodMaker) generateStatusMethodFor(die Die) {
	c.Linef("")
	c.Linef("func (d *%s) StatusDie(fn func(d %s)) %s {", die.Type, die.StatusInterface, die.Interface)
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
	c.Linef("func Test%s_MissingMethods(t *%s) {", die.Interface, c.AliasedRef("testing", "T"))
	c.Linef("	die := %s", c.AliasedRef(c.pkg.PkgPath, die.Blank))
	c.Linef("	ignore := []string{%s}", strings.Join(ignore, ", "))
	c.Linef("	diff := %s(die).Delete(ignore...)", c.AliasedRef("dies.dev/testing", "DieFieldDiff"))
	c.Linef("	if diff.Len() != 0 {")
	c.Linef("		t.Errorf(\"found missing fields for %s: %%s\", diff.List())", die.Interface)
	c.Linef("	}")
	c.Linef("}")
}
