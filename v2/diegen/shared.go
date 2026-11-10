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

// A method written by hand on a die, marked +die:shared, is shared with its
// counterpart, so it is written once. It is opt-in, as not every method of a
// die suits the other, such as one that returns an immutable die to write
// through.
//
// A method of an immutable die is shared with its mutable die: the mutable
// die is generated a method of the same signature that calls it on the die's
// [DieMutBase.DieReadOnly].
//
//	// IsTrue reports whether the condition's status is True.
//	//
//	// +die:shared
//	func (d *ConditionDie) IsTrue() bool { ... }
//
// A method of a mutable die, which writes it and returns it, is shared with
// its immutable die: the immutable die is generated a method of the same
// parameters that calls it on a mutable copy in [DieBase.DieStamp], and
// returns the immutable die of the copy in place of the mutable die, with the
// method's other results.
//
//	// +die:shared
//	func (d *ConditionDieMut) MarkTrue(reason string) (*ConditionDieMut, error) { ... }
//
// is shared as
//
//	func (d *ConditionDie) MarkTrue(reason string) (*ConditionDie, error)

// sharedMethods returns the methods marked +die:shared, by the name of the
// die they are of, immutable or mutable, in the order they are declared.
func (pg *pkgGen) sharedMethods() (map[string][]*types.Func, error) {
	dies, muts := map[string]bool{}, map[string]bool{}
	for _, tg := range pg.targets {
		dies[tg.name()+"Die"] = true
		muts[tg.name()+"DieMut"] = true
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
			case !dies[recv] && !muts[recv]:
				errs = append(errs, fmt.Sprintf("%s: +die:shared on %s.%s: only a method of a die, such as %sDie or %sDieMut, is shared", at, recv, fd.Name.Name, recv, recv))
			case !ok:
				errs = append(errs, fmt.Sprintf("%s: +die:shared on %s.%s: the method doesn't type check", at, recv, fd.Name.Name))
			case fn.Signature().TypeParams().Len() > 0:
				errs = append(errs, fmt.Sprintf("%s: +die:shared on %s.%s: a generic method isn't shared yet", at, recv, fd.Name.Name))
			case muts[recv] && pg.dieResults(fn, recv) != 1:
				errs = append(errs, fmt.Sprintf("%s: +die:shared on %s.%s: a method of a mutable die is shared when it returns the die, *%s, once", at, recv, fd.Name.Name, recv))
			default:
				shared[recv] = append(shared[recv], fn)
			}
		}
	}
	for recv, fns := range shared {
		other := recv + "Mut"
		if muts[recv] {
			other = strings.TrimSuffix(recv, "Mut")
		}
		for _, fn := range fns {
			if byHand[other+"."+fn.Name()] {
				errs = append(errs, fmt.Sprintf("%s: %s.%s is declared by hand, and is shared from %s.%s by +die:shared; remove one", pg.p.Fset.Position(fn.Pos()), other, fn.Name(), recv, fn.Name()))
			}
		}
	}
	if len(errs) > 0 {
		slices.Sort(errs)
		return nil, fmt.Errorf("%s", strings.Join(errs, "\n"))
	}
	return shared, nil
}

// dieResults returns how many of a method's results are a pointer to the
// die of the package named name.
func (pg *pkgGen) dieResults(fn *types.Func, name string) int {
	n := 0
	for i := range fn.Signature().Results().Len() {
		if pg.isDie(fn.Signature().Results().At(i).Type(), name) {
			n++
		}
	}
	return n
}

// isDie reports whether a type is a pointer to the die of this package named
// name.
func (pg *pkgGen) isDie(t types.Type, name string) bool {
	p, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return false
	}
	n, ok := types.Unalias(p.Elem()).(*types.Named)
	return ok && n.Obj().Pkg() == pg.p.Types && n.Obj().Name() == name
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

// sharedParams returns the parameters of a shared method, and the arguments
// that pass them on, and a name for each of more, such as the receiver,
// that none of them hides.
func (pg *pkgGen) sharedParams(sig *types.Signature, more ...string) (params, args, extra []string) {
	q := pg.qualifier()
	var names []string
	for i := range sig.Params().Len() {
		names = append(names, sig.Params().At(i).Name())
	}
	for _, n := range more {
		for slices.Contains(names, n) {
			n += "_"
		}
		names = append(names, n)
		extra = append(extra, n)
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
	return params, args, extra
}

// qualifier qualifies the types of other packages than this one, importing
// them.
func (pg *pkgGen) qualifier() types.Qualifier {
	return func(p *types.Package) string {
		if p.Path() == pg.p.PkgPath {
			return ""
		}
		return pg.f.use(p.Path())
	}
}

// resultList returns the results of a signature as they are declared.
func resultList(results []string) string {
	switch len(results) {
	case 0:
		return ""
	case 1:
		return " " + results[0]
	default:
		return " (" + strings.Join(results, ", ") + ")"
	}
}

// emitShared writes the mutable die's method that shares a method of the
// immutable die.
func (pg *pkgGen) emitShared(dname, mname string, fn *types.Func) {
	f := pg.f
	q := pg.qualifier()
	sig := fn.Signature()
	params, args, extra := pg.sharedParams(sig, "d")
	recv := extra[0]
	var results []string
	for i := range sig.Results().Len() {
		results = append(results, types.TypeString(sig.Results().At(i).Type(), q))
	}
	ret := ""
	if len(results) > 0 {
		ret = "return "
	}
	f.p("// %s is [%s.%s], shared by +die:shared.", fn.Name(), dname, fn.Name())
	f.p("func (%s *%s) %s(%s)%s {", recv, mname, fn.Name(), strings.Join(params, ", "), resultList(results))
	f.p("\t%s%s.DieReadOnly().%s(%s)", ret, recv, fn.Name(), strings.Join(args, ", "))
	f.p("}")
	f.p("")
}

// emitStamped writes the immutable die's method that shares a method of the
// mutable die, by calling it on a mutable copy in DieStamp: the mutable die
// it returns is the immutable die of the copy, and its other results are
// passed through.
func (pg *pkgGen) emitStamped(dname, mname string, fn *types.Func) {
	f := pg.f
	q := pg.qualifier()
	sig := fn.Signature()
	more := []string{"d", "m", "out"}
	for i := range sig.Results().Len() {
		more = append(more, fmt.Sprintf("r%d", i))
	}
	params, args, extra := pg.sharedParams(sig, more...)
	recv, m, out := extra[0], extra[1], extra[2]
	var results, vars, captures, returns []string
	for i := range sig.Results().Len() {
		t := sig.Results().At(i).Type()
		if pg.isDie(t, mname) {
			results = append(results, "*"+dname)
			captures = append(captures, "_")
			returns = append(returns, out)
			continue
		}
		r, ts := extra[3+i], types.TypeString(t, q)
		results = append(results, ts)
		vars = append(vars, r+" "+ts)
		captures = append(captures, r)
		returns = append(returns, r)
	}
	// the die returned is the copy, which DieStamp returns immutable
	call := strings.Join(captures, ", ") + " = " + fmt.Sprintf("%s.%s(%s)", m, fn.Name(), strings.Join(args, ", "))
	f.p("// %s is [%s.%s] on a copy of the die, shared by +die:shared.", fn.Name(), mname, fn.Name())
	f.p("func (%s *%s) %s(%s)%s {", recv, dname, fn.Name(), strings.Join(params, ", "), resultList(results))
	for _, v := range vars {
		f.p("\tvar %s", v)
	}
	f.p("\t%s := %s.DieStamp(func(%s *%s) {", out, recv, m, mname)
	f.p("\t\t%s", call)
	f.p("\t})")
	f.p("\treturn %s", strings.Join(returns, ", "))
	f.p("}")
	f.p("")
}
