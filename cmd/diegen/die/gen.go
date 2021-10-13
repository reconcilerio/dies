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
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"
	"strings"
	"unicode"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

var (
	dieMarker    = markers.Must(markers.MakeDefinition("die", markers.DescribesPackage, Die{}))
	fieldMarker  = markers.Must(markers.MakeDefinition("die:field", markers.DescribesPackage, Field{}))
	schemeMarker = markers.Must(markers.MakeDefinition("die:scheme", markers.DescribesPackage, Scheme{}))
)

type Scheme struct {
	Group   string `marker:"group"`
	Version string `marker:"version"`
}

type Die struct {
	Target       string   `marker:"target"`
	Object       bool     `marker:"object,optional"`
	IgnoreFields []string `marker:"ignore,optional"`

	Name          string `marker:",optional"`
	Type          string `marker:",optional"`
	TargetPackage string `marker:",optional"`
	TargetType    string `marker:",optional"`
	Blank         string `marker:",optional"`
	SpecName      string `marker:",optional"`
	SpecBlank     string `marker:",optional"`
	SpecType      string `marker:",optional"`
	StatusName    string `marker:",optional"`
	StatusBlank   string `marker:",optional"`
	StatusType    string `marker:",optional"`
}

func (d *Die) Default() {
	i := strings.LastIndex(d.Target, ".")
	d.TargetPackage = d.Target[0:i]
	d.TargetType = d.Target[i+1:]
	if d.Name == "" {
		d.Name = d.TargetType
	}
	if d.Type == "" {
		d.Type = fmt.Sprintf("%sDie", d.Name)
	}
	if d.Blank == "" {
		d.Blank = fmt.Sprintf("%sBlank", d.Name)
	}
	if d.SpecName == "" {
		d.SpecName = fmt.Sprintf("%sSpec", d.Name)
	}
	if d.SpecBlank == "" {
		d.SpecBlank = fmt.Sprintf("%sBlank", d.SpecName)
	}
	if d.SpecType == "" {
		d.SpecType = fmt.Sprintf("%sDie", d.SpecName)
	}
	if d.StatusName == "" {
		d.StatusName = fmt.Sprintf("%sStatus", d.Name)
	}
	if d.StatusBlank == "" {
		d.StatusBlank = fmt.Sprintf("%sBlank", d.StatusName)
	}
	if d.StatusType == "" {
		d.StatusType = fmt.Sprintf("%sDie", d.StatusName)
	}
	if d.IgnoreFields == nil {
		d.IgnoreFields = []string{}
	}
	if d.Object {
		d.IgnoreFields = append(d.IgnoreFields, "TypeMeta", "ObjectMeta")
	}
}

type Field struct {
	Receiver string `marker:"receiver"`
	Name     string `marker:"name"`
	Type     string `marker:"type"`

	TypePrefix  string `marker:",optional"`
	TypePackage string `marker:",optional"`
}

func (d *Field) Default() {
	if i := strings.IndexFunc(d.Type, unicode.IsLetter); i >= 0 {
		d.TypePrefix = d.Type[0:i]
		d.Type = d.Type[i:]

		// spread slices
		d.TypePrefix = strings.Replace(d.TypePrefix, "[]", "...", 1)
	}
	if i := strings.LastIndex(d.Type, "."); i >= 0 {
		d.TypePackage = d.Type[0:i]
		d.Type = d.Type[i+1:]
		if strings.HasPrefix(d.TypePackage, "*") {
			d.TypePackage = d.TypePackage[1:]
			d.Type = "*" + d.Type
		}
	}
}

type Generator struct {
	// HeaderFile specifies the header text (e.g. license) to prepend to generated files.
	HeaderFile string `marker:",optional"`
	// Year specifies the year to substitute for " YEAR" in the header file.
	Year string `marker:",optional"`
}

func (Generator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool {
		// ignore interfaces
		_, isIface := node.(*ast.InterfaceType)
		return !isIface
	}
}

func (Generator) RegisterMarkers(into *markers.Registry) error {
	if err := into.Register(dieMarker); err != nil {
		return err
	}
	into.AddHelp(dieMarker, markers.SimpleHelp("die", "generates a die for the type"))
	if err := into.Register(fieldMarker); err != nil {
		return err
	}
	into.AddHelp(fieldMarker, markers.SimpleHelp("die:field", "generates a field mutator for the die"))
	if err := into.Register(schemeMarker); err != nil {
		return err
	}
	into.AddHelp(schemeMarker, markers.SimpleHelp("die:scheme", "generates a scheme for the package that object dies are registered with"))

	return nil
}

func (d Generator) Generate(ctx *genall.GenerationContext) error {
	var headerText string

	if d.HeaderFile != "" {
		headerBytes, err := ctx.ReadFile(d.HeaderFile)
		if err != nil {
			return err
		}
		headerText = string(headerBytes)
	}
	headerText = strings.ReplaceAll(headerText, " YEAR", " "+d.Year)

	objGenCtx := ObjectGenCtx{
		Collector:  ctx.Collector,
		Checker:    ctx.Checker,
		HeaderText: headerText,
	}

	for _, root := range ctx.Roots {
		outContents, testContents := objGenCtx.generateForPackage(root)
		if outContents != nil {
			writeOut(ctx, root, "zz_generated.die.go", outContents)
		}
		if testContents != nil {
			writeOut(ctx, root, "zz_generated.die_test.go", testContents)
		}
	}

	return nil
}

// ObjectGenCtx contains the common info for generating deepcopy implementations.
// It mostly exists so that generating for a package can be easily tested without
// requiring a full set of output rules, etc.
type ObjectGenCtx struct {
	Collector  *markers.Collector
	Checker    *loader.TypeChecker
	HeaderText string
}

// writeHeader writes out the build tag, package declaration, and imports
func writeHeader(pkg *loader.Package, out io.Writer, packageName string, imports *importsList, headerText string) {
	// NB(directxman12): blank line after build tags to distinguish them from comments
	_, err := fmt.Fprintf(out, `// +build !ignore_autogenerated

%[3]s

// Code generated by diegen. DO NOT EDIT.

package %[1]s

import (
%[2]s
)

`, packageName, strings.Join(imports.ImportSpecs(), "\n"), headerText)
	if err != nil {
		pkg.AddError(err)
	}

}

// generateForPackage generates DeepCopy and runtime.Object implementations for
// types in the given package, writing the formatted result to given writer.
// May return nil if source could not be generated.
func (ctx *ObjectGenCtx) generateForPackage(root *loader.Package) ([]byte, []byte) {
	ctx.Checker.Check(root)

	root.NeedTypesInfo()

	imports := &importsList{
		byPath:  make(map[string]string),
		byAlias: make(map[string]string),
		pkg:     root,
	}
	// avoid confusing aliases by "reserving" the root package's name as an alias
	imports.byAlias[root.Name] = ""

	testImports := &importsList{
		byPath:  make(map[string]string),
		byAlias: make(map[string]string),
		pkg:     root,
	}
	// avoid confusing aliases by "reserving" the root package's name as an alias
	testImports.byAlias[root.Name] = ""

	markerSet, err := markers.PackageMarkers(ctx.Collector, root)
	if err != nil {
		root.AddError(err)
	}

	outContent := new(bytes.Buffer)
	testContent := new(bytes.Buffer)

	dies := []Die{}
	dieSet := sets.NewString()
	for _, markerValue := range markerSet[dieMarker.Name] {
		die := markerValue.(Die)
		die.Default()
		dies = append(dies, die)
		dieSet.Insert(die.Name)
	}
	fieldMap := map[string][]Field{}
	for _, markerValue := range markerSet[fieldMarker.Name] {
		field := markerValue.(Field)
		field.Default()
		if _, ok := fieldMap[field.Receiver]; !ok {
			fieldMap[field.Receiver] = []Field{}
		}
		fieldMap[field.Receiver] = append(fieldMap[field.Receiver], field)
	}

	copyCtx := &copyMethodMaker{
		pkg:         root,
		importsList: imports,
		codeWriter:  &codeWriter{out: outContent},
		test: &copyMethodMaker{
			pkg:         root,
			importsList: testImports,
			codeWriter:  &codeWriter{out: testContent},
			dies:        dieSet,
		},
		dies: dieSet,
	}

	schemeValues := markerSet[schemeMarker.Name]
	if len(schemeValues) == 1 {
		scheme := schemeValues[0].(Scheme)
		copyCtx.GenerateSchemeFor(scheme)
	}

	for _, die := range dies {
		fmt.Printf("Generating die for %q\n", die.Name)
		copyCtx.GenerateMethodsFor(die)

		// print fields for this die
		for _, field := range fieldMap[die.Type] {
			fmt.Printf("Generating field %q for %q\n", field.Name, die.Name)
			copyCtx.GenerateFieldFor(field, die)
		}
		delete(fieldMap, die.Type)
	}

	// print fields not generated for a known die
	for _, fields := range fieldMap {
		for _, field := range fields {
			fmt.Printf("Skipping field %q for unknown die %q\n", field.Name, field.Receiver)
		}
	}

	outBytes := ctx.outputBytes(root, imports, outContent.Bytes())
	testBytes := ctx.outputBytes(root, testImports, testContent.Bytes())

	return outBytes, testBytes
}

func (ctx *ObjectGenCtx) outputBytes(root *loader.Package, imports *importsList, content []byte) []byte {
	if len(content) == 0 {
		return nil
	}

	outContentWithHeader := new(bytes.Buffer)
	writeHeader(root, outContentWithHeader, root.Name, imports, ctx.HeaderText)
	outContentWithHeader.Write(content)

	outBytes := outContentWithHeader.Bytes()
	formattedBytes, err := format.Source(outBytes)
	if err != nil {
		root.AddError(err)
		// we still write the invalid source to disk to figure out what went wrong
	} else {
		outBytes = formattedBytes
	}

	return outBytes
}

// writeFormatted outputs the given code, after gofmt-ing it.  If we couldn't gofmt,
// we write the unformatted code for debugging purposes.
func writeOut(ctx *genall.GenerationContext, root *loader.Package, filename string, outBytes []byte) {
	outputFile, err := ctx.Open(root, filename)
	if err != nil {
		root.AddError(err)
		return
	}
	defer outputFile.Close()
	n, err := outputFile.Write(outBytes)
	if err != nil {
		root.AddError(err)
		return
	}
	if n < len(outBytes) {
		root.AddError(io.ErrShortWrite)
	}
}
