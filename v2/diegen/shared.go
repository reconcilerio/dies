/*
Copyright 2026 the original author or authors.

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

package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"slices"
	"strings"
)

// A method written by hand on an immutable die, marked +die:shared, is shared
// with its mutable die: the mutable die is generated a method of the same
// signature that calls it on the die's [DieMutBase.DieReadOnly], so it is
// written once. It is opt-in, as not every method of an immutable die suits a
// mutable one, such as one that returns an immutable die to write through.
//
//	// IsTrue reports whether the condition's status is True.
//	//
//	// +die:shared
//	func (d *ConditionDie) IsTrue() bool { ... }

// sharedMethods returns the methods marked +die:shared, by the name of the
// die they are of, in the order they are declared.
func (pg *pkgGen) sharedMethods() (map[string][]*types.Func, error) {
	dies := map[string]bool{}
	for _, tg := range pg.targets {
		dies[tg.name()+"Die"] = true
	}
	shared := map[string][]*types.Func{}
	byHand := map[string]bool{}
	var errs []string
	for _, f := range pg.p.Syntax {
		switch filepath.Base(pg.p.Fset.File(f.Pos()).Name()) {
		case typesFile, diesFile, acFile:
			continue
		}
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv == nil || len(fd.Recv.List) != 1 {
				continue
			}
			recv := receiverName(fd.Recv.List[0].Type)
			byHand[recv+"."+fd.Name.Name] = true
			if _, ok := findDieMarker(commentLines(fd.Doc), "die:shared"); !ok {
				continue
			}
			at := pg.p.Fset.Position(fd.Pos())
			fn, ok := pg.p.TypesInfo.Defs[fd.Name].(*types.Func)
			switch {
			case !dies[recv]:
				errs = append(errs, fmt.Sprintf("%s: +die:shared on %s.%s: only a method of an immutable die, such as %sDie, is shared with its mutable die", at, recv, fd.Name.Name, strings.TrimSuffix(recv, "DieMut")))
			case !ok:
				errs = append(errs, fmt.Sprintf("%s: +die:shared on %s.%s: the method doesn't type check", at, recv, fd.Name.Name))
			case fn.Signature().TypeParams().Len() > 0:
				errs = append(errs, fmt.Sprintf("%s: +die:shared on %s.%s: a generic method isn't shared yet", at, recv, fd.Name.Name))
			default:
				shared[recv] = append(shared[recv], fn)
			}
		}
	}
	for recv, fns := range shared {
		for _, fn := range fns {
			if byHand[recv+"Mut."+fn.Name()] {
				errs = append(errs, fmt.Sprintf("%s: %sMut.%s is declared by hand, and is shared from %s.%s by +die:shared; remove one", pg.p.Fset.Position(fn.Pos()), recv, fn.Name(), recv, fn.Name()))
			}
		}
	}
	if len(errs) > 0 {
		slices.Sort(errs)
		return nil, fmt.Errorf("%s", strings.Join(errs, "\n"))
	}
	return shared, nil
}

// receiverName returns the name of a receiver's type, without its pointer.
func receiverName(e ast.Expr) string {
	if s, ok := e.(*ast.StarExpr); ok {
		e = s.X
	}
	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// emitShared writes the mutable die's method that shares a method of the
// immutable die.
func (pg *pkgGen) emitShared(dname, mname string, fn *types.Func) {
	f := pg.f
	q := func(p *types.Package) string {
		if p.Path() == pg.p.PkgPath {
			return ""
		}
		return f.use(p.Path())
	}
	sig := fn.Signature()
	var params, args, names []string
	for i := range sig.Params().Len() {
		names = append(names, sig.Params().At(i).Name())
	}
	recv := "d"
	for slices.Contains(names, recv) {
		recv += "_"
	}
	for i := range sig.Params().Len() {
		v := sig.Params().At(i)
		name := v.Name()
		if name == "" || name == "_" {
			name = fmt.Sprintf("p%d", i)
		}
		t, arg := types.TypeString(v.Type(), q), name
		if sig.Variadic() && i == sig.Params().Len()-1 {
			t, arg = "..."+types.TypeString(v.Type().(*types.Slice).Elem(), q), name+"..."
		}
		params = append(params, name+" "+t)
		args = append(args, arg)
	}
	var results []string
	for i := range sig.Results().Len() {
		results = append(results, types.TypeString(sig.Results().At(i).Type(), q))
	}
	result, ret := "", ""
	switch len(results) {
	case 0:
	case 1:
		result, ret = " "+results[0], "return "
	default:
		result, ret = " ("+strings.Join(results, ", ")+")", "return "
	}
	f.p("// %s is [%s.%s], shared by +die:shared.", fn.Name(), dname, fn.Name())
	f.p("func (%s *%s) %s(%s)%s {", recv, mname, fn.Name(), strings.Join(params, ", "), result)
	f.p("\t%s%s.DieReadOnly().%s(%s)", ret, recv, fn.Name(), strings.Join(args, ", "))
	f.p("}")
	f.p("")
}
