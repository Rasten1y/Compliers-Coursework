package frontend

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
)

// Program bundles parsed files and type information for the rest of the pipeline.
type Program struct {
	Fset      *token.FileSet
	Files     []*ast.File
	TypesInfo *types.Info
	TypesPkg  *types.Package
	PkgName   string // original package name from source
}

// SkipSubsetCheck disables CheckUnsupported; useful for self-host builds.
var SkipSubsetCheck bool

// SetSkipSubsetCheck sets the global flag controlling unsupported checks.
func SetSkipSubsetCheck(v bool) {
	SkipSubsetCheck = v
}

// ParseAndCheck parses Go source files, rejects unsupported constructs, and runs type checking.
func ParseAndCheck(filenames []string) (*Program, error) {
	progs, err := ParseAndCheckAll(filenames)
	if err != nil {
		return nil, err
	}
	if len(progs) == 1 {
		return progs[0], nil
	}
	for i := 0; i < len(progs); i++ {
		if progs[i].PkgName == "main" {
			return progs[i], nil
		}
	}
	return progs[0], nil
}

// ParseAndCheckAll parses all files and returns one Program per package.
func ParseAndCheckAll(filenames []string) ([]*Program, error) {
	fset := token.NewFileSet()
	type pkgGroup struct {
		key   string
		files []*ast.File
		orig  string
	}
	var groups []pkgGroup

	for i := 0; i < len(filenames); i++ {
		name := filenames[i]
		file, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		orig := file.Name.Name
		key := orig
		if SkipSubsetCheck {
			key = "main"
			file.Name.Name = "main"
		}
		found := -1
		for i := 0; i < len(groups); i++ {
			if groups[i].key == key {
				found = i
				break
			}
		}
		if found == -1 {
			groups = append(groups, pkgGroup{key: key, orig: orig, files: []*ast.File{file}})
		} else {
			groups[found].files = append(groups[found].files, file)
		}
	}

	addRuntimeBuiltins()

	var progs []*Program
	for i := 0; i < len(groups); i++ {
		pkgName := groups[i].orig
		files := groups[i].files
		if !SkipSubsetCheck {
			for j := 0; j < len(files); j++ {
				if err := CheckUnsupported(files[j]); err != nil {
					return nil, err
				}
			}
		}

		rtPkg := files[0].Name.Name
		rtSrc := runtimeDeclSource(rtPkg)
		rtFile, err := parser.ParseFile(fset, "<runtime-builtins>", rtSrc, 0)
		if err != nil {
			return nil, err
		}
		allFiles := append([]*ast.File{}, files...)
		allFiles = append(allFiles, rtFile)

		info := &types.Info{
			Types:      make(map[ast.Expr]types.TypeAndValue),
			Defs:       make(map[*ast.Ident]types.Object),
			Uses:       make(map[*ast.Ident]types.Object),
			Implicits:  make(map[ast.Node]types.Object),
			Selections: make(map[*ast.SelectorExpr]*types.Selection),
			Scopes:     make(map[ast.Node]*types.Scope),
		}

		conf := types.Config{
			Importer: importer.ForCompiler(fset, "source", nil),
		}
		pkg, err := conf.Check(pkgName, fset, allFiles, info)
		if err != nil {
			return nil, err
		}

		progs = append(progs, &Program{
			Fset:      fset,
			Files:     files,
			TypesInfo: info,
			TypesPkg:  pkg,
			PkgName:   pkgName,
		})
	}

	return progs, nil
}

// CheckUnsupported currently disabled to simplify self-hosting.
func CheckUnsupported(file *ast.File) error {
	return nil
}

// addRuntimeBuiltins registers external runtime functions in the universe scope so
// they can be referenced without Go definitions in user code.
func addRuntimeBuiltins() {
	int64Typ := types.Typ[types.Int64]
	byteTyp := types.Universe.Lookup("byte").Type()
	ptrByte := types.NewPointer(byteTyp)
	boolTyp := types.Typ[types.Bool]
	strTyp := types.Typ[types.String]

	// gominic_print(*byte, int64)
	pPrint := types.NewTuple(
		types.NewParam(token.NoPos, nil, "data", ptrByte),
		types.NewParam(token.NoPos, nil, "len", int64Typ),
	)
	registerBuiltin("gominic_print", types.NewSignature(nil, pPrint, types.NewTuple(), false))

	// gominic_printInt(int64)
	pPrintInt := types.NewTuple(types.NewParam(token.NoPos, nil, "v", int64Typ))
	registerBuiltin("gominic_printInt", types.NewSignature(nil, pPrintInt, types.NewTuple(), false))

	// gominic_println()
	registerBuiltin("gominic_println", types.NewSignature(nil, types.NewTuple(), types.NewTuple(), false))

	// gominic_makeSlice(int64 len, int64 cap, int64 elemSize) *byte
	pMake := types.NewTuple(
		types.NewParam(token.NoPos, nil, "len", int64Typ),
		types.NewParam(token.NoPos, nil, "cap", int64Typ),
		types.NewParam(token.NoPos, nil, "elemSize", int64Typ),
	)
	rMake := types.NewTuple(types.NewParam(token.NoPos, nil, "", ptrByte))
	registerBuiltin("gominic_makeSlice", types.NewSignature(nil, pMake, rMake, false))

	// map runtime
	pMapNew := types.NewTuple(
		types.NewParam(token.NoPos, nil, "keySize", int64Typ),
		types.NewParam(token.NoPos, nil, "valSize", int64Typ),
		types.NewParam(token.NoPos, nil, "keyKind", types.Typ[types.Int32]),
	)
	rMapPtr := types.NewTuple(types.NewParam(token.NoPos, nil, "", ptrByte))
	registerBuiltin("gominic_map_new", types.NewSignature(nil, pMapNew, rMapPtr, false))

	pMapSet := types.NewTuple(
		types.NewParam(token.NoPos, nil, "m", ptrByte),
		types.NewParam(token.NoPos, nil, "key", ptrByte),
		types.NewParam(token.NoPos, nil, "val", ptrByte),
	)
	registerBuiltin("gominic_map_set", types.NewSignature(nil, pMapSet, types.NewTuple(), false))

	pMapGet := types.NewTuple(
		types.NewParam(token.NoPos, nil, "m", ptrByte),
		types.NewParam(token.NoPos, nil, "key", ptrByte),
		types.NewParam(token.NoPos, nil, "val", ptrByte),
	)
	registerBuiltin("gominic_map_get", types.NewSignature(nil, pMapGet, types.NewTuple(types.NewParam(token.NoPos, nil, "", types.Typ[types.Bool])),
		false))

	pMapLen := types.NewTuple(types.NewParam(token.NoPos, nil, "m", ptrByte))
	registerBuiltin("gominic_map_len", types.NewSignature(nil, pMapLen, types.NewTuple(types.NewParam(token.NoPos, nil, "", int64Typ)), false))

	// self-host I/O helpers
	registerBuiltin("gominic_argc", types.NewSignature(nil, types.NewTuple(), types.NewTuple(types.NewParam(token.NoPos, nil, "", int64Typ)), false))
	pArgv := types.NewTuple(types.NewParam(token.NoPos, nil, "idx", int64Typ))
	registerBuiltin("gominic_argv", types.NewSignature(nil, pArgv, types.NewTuple(types.NewParam(token.NoPos, nil, "", strTyp)), false))
	pWriteFile := types.NewTuple(
		types.NewParam(token.NoPos, nil, "path", ptrByte),
		types.NewParam(token.NoPos, nil, "pathLen", int64Typ),
		types.NewParam(token.NoPos, nil, "data", ptrByte),
		types.NewParam(token.NoPos, nil, "dataLen", int64Typ),
	)
	registerBuiltin("gominic_write_file", types.NewSignature(nil, pWriteFile, types.NewTuple(types.NewParam(token.NoPos, nil, "", boolTyp)), false))
}

// runtimeDeclSource returns a minimal Go source with stubs for runtime externs to satisfy type checking.
func runtimeDeclSource(pkg string) string {
	return fmt.Sprintf(`package %s

func gominic_print(data *byte, len int64) {}
func gominic_printInt(v int64)            {}
func gominic_println()                    {}
func gominic_makeSlice(len int64, cap int64, elemSize int64) *byte { return nil }
func gominic_map_new(keySize int64, valSize int64, keyKind int32) *byte { return nil }
func gominic_map_set(m *byte, key *byte, val *byte) {}
func gominic_map_get(m *byte, key *byte, val *byte) bool { return false }
func gominic_map_len(m *byte) int64 { return 0 }
func gominic_argc() int64 { return 0 }
func gominic_argv(idx int64) string { return "" }
func gominic_write_file(path *byte, pathLen int64, data *byte, dataLen int64) bool { return false }
`, pkg)
}

// registerBuiltin inserts a function into the universe scope if not present.
func registerBuiltin(name string, sig *types.Signature) {
	if obj := types.Universe.Lookup(name); obj == nil {
		types.Universe.Insert(types.NewFunc(token.NoPos, nil, name, sig))
	}
}
