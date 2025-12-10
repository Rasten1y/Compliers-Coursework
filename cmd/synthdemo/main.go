package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gominic/backend"
	"gominic/ir"
)

// synthdemo: минимальный драйвер стадии синтеза без frontend.
// Он берет захардкоженный пример (эквивалент простому AST) и генерирует LLVM IR.

func main() {
	out := flag.String("o", "", "файл для вывода LLVM IR (по умолчанию stdout)")
	flag.Parse()

	mod := buildSampleModule()
	text := backend.EmitModule(mod)

	if *out == "" {
		fmt.Print(text)
		return
	}
	if err := os.WriteFile(*out, []byte(text), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write failed: %v\n", err)
		os.Exit(1)
	}
}

// buildSampleModule создаёт модуль с одной функцией main:
// print("hello from synthdemo"), printInt(42), println().
func buildSampleModule() *ir.Module {
	mod := &ir.Module{
		TargetTriple: "x86_64-pc-windows-msvc",
		DataLayout:   "e-m:w-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128",
	}

	helloPtr, helloLen := addStringGlobal(mod, ".str.hello", "hello from synthdemo")

	mainFn := &ir.Function{
		Name:    "main",
		Params:  nil,
		Results: nil,
	}
	entry := ir.NewBasicBlock("entry")

	// call void @gominic_print(i8* ptr, i64 len)
	entry.Append(ir.Call{
		Name: "gominic_print",
		Args: []ir.Value{helloPtr, helloLen},
		Ret:  ir.Void,
	})

	// call void @gominic_printInt(i64 42)
	entry.Append(ir.Call{
		Name: "gominic_printInt",
		Args: []ir.Value{ir.NewConstant("42", ir.I64)},
		Ret:  ir.Void,
	})

	// call void @gominic_println()
	entry.Append(ir.Call{
		Name: "gominic_println",
		Args: nil,
		Ret:  ir.Void,
	})

	entry.Terminator = ir.Return{}
	mainFn.Blocks = append(mainFn.Blocks, entry)
	mod.AddFunction(mainFn)
	return mod
}

// addStringGlobal создаёт приватную строку в разделе Globals и возвращает указатель/длину.
func addStringGlobal(mod *ir.Module, name, content string) (ptr ir.Value, length ir.Value) {
	bytes := []byte(content)
	enc := escapeCString(bytes)
	arr := ir.ArrayType{Len: len(bytes) + 1, Elem: ir.I8}
	mod.Globals = append(mod.Globals, ir.Global{
		Name:    name,
		Type:    arr,
		Value:   fmt.Sprintf("c\"%s\"", enc),
		Align:   1,
		Private: true,
	})
	gep := fmt.Sprintf("getelementptr inbounds (%s, %s* @%s, i32 0, i32 0)", arr.String(), arr.String(), name)
	ptr = ir.NewConstant(gep, ir.PtrTo(ir.I8))
	length = ir.NewConstant(strconv.Itoa(len(bytes)), ir.I64)
	return
}

func escapeCString(b []byte) string {
	var sb strings.Builder
	for _, ch := range b {
		fmt.Fprintf(&sb, "\\%02X", ch)
	}
	sb.WriteString("\\00")
	return sb.String()
}
