package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"gominic/backend"
	"gominic/frontend"
	"gominic/ir"
)

func main() {
	var output string
	var emitIR bool
	var stopObj bool
	var verbose bool
	var llOut string
	var clangPath string
	var skipCheck bool
	var demo bool

	flag.StringVar(&output, "o", "a.out", "output file")
	flag.BoolVar(&emitIR, "S", false, "emit LLVM IR and exit")
	flag.BoolVar(&stopObj, "c", false, "emit object file and exit")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&llOut, "ll", "", "write LLVM IR to file instead of stdout when using -S")
	flag.StringVar(&clangPath, "cc", "clang", "C compiler (used also to compile .ll)")
	flag.BoolVar(&skipCheck, "skip-check", false, "skip subset checks (useful for self-hosting)")
	flag.BoolVar(&demo, "demo", false, "use hardcoded demo IR instead of parsing Go files")
	flag.Parse()

	files := flag.Args()
	var mod *ir.Module
	if demo {
		if verbose {
			vprintf("using demo IR (hardcoded)\n")
		}
		mod = buildDemoModule()
	} else {
		if len(files) == 0 {
			flag.Usage()
			exit(2)
		}
		if verbose {
			vprintf("compiling files\n")
		}
		frontend.SetSkipSubsetCheck(skipCheck)
		progs, err := frontend.ParseAndCheckAll(files)
		if err != nil {
			fatalf("frontend failed: " + err.Error())
		}
		if verbose {
			for i := 0; i < len(progs); i++ {
				p := progs[i]
				vprintf("type checked package: " + p.TypesPkg.Name() + "\n")
			}
		}
		mod, err = frontend.BuildModule(progs)
		if err != nil {
			fatalf("IR lowering failed: " + err.Error())
		}
	}

	irText := backend.EmitModule(mod)

	// runtime objects selection (demo нуждается только в print.c)
	rtEntries := []struct {
		src string
		out string
	}{
		{"runtime/print.c", "runtime/print.o"},
	}
	if !demo {
		rtEntries = append(rtEntries, []struct {
			src string
			out string
		}{
			{"runtime/map.c", "runtime/map.o"},
			{"runtime/io.c", "runtime/io.o"},
		}...)
	}

	if emitIR {
		if llOut != "" {
			if err := os.WriteFile(llOut, []byte(irText), 0644); err != nil {
				fatalf("write ll failed: " + err.Error())
			}
		} else {
			os.Stdout.WriteString(irText)
		}
	} else if stopObj {
		llPath := llOut
		if llPath == "" {
			llPath = deriveLLPath(output)
		}
		if err := os.WriteFile(llPath, []byte(irText), 0644); err != nil {
			fatalf("write ll failed: " + err.Error())
		}
		objPath := output
		if filepath.Ext(objPath) != ".o" {
			objPath += ".o"
		}
		if err := compileLLWithClang(clangPath, llPath, objPath, verbose); err != nil {
			vprintf("compile ll failed: " + err.Error() + "\n")
			vprintf("You can try manually: " + clangPath + " -c -o " + objPath + " " + llPath + "\n")
			exit(1)
		}
		if verbose {
			vprintf("wrote object " + objPath + "\n")
		}
	} else {
		llPath := llOut
		if llPath == "" {
			llPath = deriveLLPath(output)
		}
		if err := os.WriteFile(llPath, []byte(irText), 0644); err != nil {
			fatalf("write ll failed: " + err.Error())
		}
		objPath := output
		if filepath.Ext(objPath) != ".o" {
			objPath = output + ".o"
		}
		if err := compileLLWithClang(clangPath, llPath, objPath, verbose); err != nil {
			vprintf("compile ll failed: " + err.Error() + "\n")
			vprintf("You can try manually: " + clangPath + " -c -o " + objPath + " " + llPath + "\n")
			exit(1)
		}
		if err := buildRuntimeObjs(clangPath, rtEntries, verbose); err != nil {
			vprintf("building runtime obj failed: " + err.Error() + "\n")
			exit(1)
		}
		if err := linkExecutable(clangPath, output, objPath, rtEntries, verbose); err != nil {
			vprintf("link failed: " + err.Error() + "\n")
			vprintf("You can try manually: " + clangPath + " -o " + output + " " + objPath + " runtime/print.o runtime/map.o runtime/io.o\n")
			exit(1)
		}
		if verbose {
			vprintf("built executable " + output + "\n")
		}
	}
}

func deriveLLPath(output string) string {
	ext := filepath.Ext(output)
	base := output
	if ext != "" {
		base = output[:len(output)-len(ext)]
	}
	return base + ".ll"
}

func compileLLWithClang(cc, llPath, objPath string, verbose bool) error {
	args := []string{"-c", "-o", objPath, llPath}
	if verbose {
		vprintf("running clang to compile ll\n")
	}
	cmd := exec.Command(cc, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New("clang compile ll: " + err.Error() + "\n" + string(out))
	}
	return nil
}

func buildRuntimeObjs(cc string, entries []struct {
	src string
	out string
}, verbose bool) error {
	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		args := []string{"-c", "-o", entry.out, entry.src}
		if verbose {
			vprintf("running clang to build runtime obj\n")
		}
		cmd := exec.Command(cc, args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return errors.New("cc runtime: " + err.Error() + "\n" + string(out))
		}
	}
	return nil
}

func linkExecutable(cc, output, objPath string, rtEntries []struct {
	src string
	out string
}, verbose bool) error {
	args := []string{"-o", output, objPath}
	for i := 0; i < len(rtEntries); i++ {
		args = append(args, rtEntries[i].out)
	}
	if verbose {
		vprintf("running clang to link\n")
	}
	cmd := exec.Command(cc, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New("link: " + err.Error() + "\n" + string(out))
	}
	return nil
}

func fatalf(msg string) {
	os.Stderr.WriteString(msg + "\n")
	exit(1)
}

func vprintf(msg string) {
	os.Stderr.WriteString(msg)
}

func exit(code int) {
	os.Exit(code)
}

// buildDemoModule создаёт модуль с одной функцией main, вызывающей рантайм-принты.
func buildDemoModule() *ir.Module {
	mod := &ir.Module{
		TargetTriple: "x86_64-pc-windows-msvc",
		DataLayout:   "e-m:w-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128",
	}

	helloPtr, helloLen := addStringGlobal(mod, ".str.demo", "hello from demo")

	mainFn := &ir.Function{
		Name:    "main",
		Params:  nil,
		Results: nil,
	}
	entry := ir.NewBasicBlock("entry")
	entry.Append(ir.Call{Name: "gominic_print", Args: []ir.Value{helloPtr, helloLen}, Ret: ir.Void})
	entry.Append(ir.Call{Name: "gominic_printInt", Args: []ir.Value{ir.NewConstant("42", ir.I64)}, Ret: ir.Void})
	entry.Append(ir.Call{Name: "gominic_println", Args: nil, Ret: ir.Void})
	entry.Terminator = ir.Return{}
	mainFn.Blocks = append(mainFn.Blocks, entry)
	mod.AddFunction(mainFn)
	return mod
}

// addStringGlobal кладёт приватную строку в Globals и возвращает (ptr,len).
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
