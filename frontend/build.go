package frontend

import (
	"go/ast"
	"runtime"

	"gominic/ir"
)

// BuildModule lowers all packages into a single IR module.
func BuildModule(progs []*Program) (*ir.Module, error) {
	triple, layout := defaultTarget()
	mod := &ir.Module{
		TargetTriple: triple,
		DataLayout:   layout,
	}
	for i := 0; i < len(progs); i++ {
		if err := buildInto(progs[i], mod); err != nil {
			return nil, err
		}
	}
	return mod, nil
}

// BuildIR builds a module for a single package (legacy helper).
func BuildIR(prog *Program) (*ir.Module, error) {
	triple, layout := defaultTarget()
	mod := &ir.Module{
		TargetTriple: triple,
		DataLayout:   layout,
	}
	if err := buildInto(prog, mod); err != nil {
		return nil, err
	}
	return mod, nil
}

func defaultTarget() (string, string) {
	if runtime.GOOS == "windows" {
		// MSVC-style triple; datalayout matches clang-cl defaults.
		return "x86_64-pc-windows-msvc", "e-m:w-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
	}
	// Fallback to Linux defaults.
	return "x86_64-pc-linux-gnu", "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
}

// buildInto lowers a single package into an existing module.
func buildInto(prog *Program, mod *ir.Module) error {
	l := &lowerer{
		prog: prog,
		mod:  mod,
	}
	if prog.PkgName != "" && prog.PkgName != "main" {
		l.prefix = prog.PkgName + "."
	}

	for i := 0; i < len(prog.Files); i++ {
		if err := l.lowerGlobals(prog.Files[i].Decls); err != nil {
			return err
		}
	}

	for i := 0; i < len(prog.Files); i++ {
		file := prog.Files[i]
		for j := 0; j < len(file.Decls); j++ {
			decl := file.Decls[j]
			fnDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if err := l.lowerFunc(fnDecl); err != nil {
				return err
			}
		}
	}
	return nil
}
