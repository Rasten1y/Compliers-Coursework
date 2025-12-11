package main

import (
	"errors"
	"flag"
	"os"
	"os/exec"

	"gominic/backend"
	"gominic/frontend"
	"gominic/ir"
)

type parseResult struct {
	progs []*frontend.Program
	err   error
}

func parseAndCheckAllWrapper(files []string) parseResult {
	var result parseResult
	result.progs, result.err = frontend.ParseAndCheckAll(files)
	return result
}

type buildModuleResult struct {
	mod *ir.Module
	err error
}

func buildModuleWrapper(progs []*frontend.Program) buildModuleResult {
	var result buildModuleResult
	result.mod, result.err = frontend.BuildModule(progs)
	return result
}

type combinedOutputResult struct {
	out []byte
	err error
}

func combinedOutputWrapper(cmd *exec.Cmd) combinedOutputResult {
	var result combinedOutputResult
	result.out, result.err = cmd.CombinedOutput()
	return result
}

func main() {
	var output string
	var emitIR bool
	var stopObj bool
	var verbose bool
	var llOut string
	var clangPath string
	var skipCheck bool

	flag.StringVar(&output, "o", "a.out", "output file")
	flag.BoolVar(&emitIR, "S", false, "emit LLVM IR and exit")
	flag.BoolVar(&stopObj, "c", false, "emit object file and exit")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&llOut, "ll", "", "write LLVM IR to file instead of stdout when using -S")
	flag.StringVar(&clangPath, "cc", "clang", "C compiler (used also to compile .ll)")
	flag.BoolVar(&skipCheck, "skip-check", false, "skip subset checks (useful for self-hosting)")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		flag.Usage()
		exit(2)
	}
	if verbose {
		vprintf("compiling files\n")
	}
	frontend.SetSkipSubsetCheck(skipCheck)
	parseRes := parseAndCheckAllWrapper(files)
	if parseRes.err != nil {
		fatalf("frontend failed: " + parseRes.err.Error())
	}
	progs := parseRes.progs
	if verbose {
		for i := 0; i < len(progs); i = i + 1 {
			p := progs[i]
			vprintf("type checked package: " + p.TypesPkg.Name() + "\n")
		}
	}
	buildRes := buildModuleWrapper(progs)
	if buildRes.err != nil {
		fatalf("IR lowering failed: " + buildRes.err.Error())
	}
	mod := buildRes.mod

	irText := backend.EmitModule(mod)

	// runtime objects selection
	rtEntries := []struct {
		src string
		out string
	}{
		{"runtime/print.c", "runtime/print.o"},
		{"runtime/map.c", "runtime/map.o"},
		{"runtime/io.c", "runtime/io.o"},
	}

	if emitIR {
		if llOut != "" {
			err := os.WriteFile(llOut, []byte(irText), 0644)
			if err != nil {
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
		err := os.WriteFile(llPath, []byte(irText), 0644)
		if err != nil {
			fatalf("write ll failed: " + err.Error())
		}
		objPath := output
		extO := ".o"
		if fileExt(objPath) != extO {
			objPath = objPath + extO
		}
		err = compileLLWithClang(clangPath, llPath, objPath, verbose)
		if err != nil {
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
		err := os.WriteFile(llPath, []byte(irText), 0644)
		if err != nil {
			fatalf("write ll failed: " + err.Error())
		}
		objPath := output
		extO := ".o"
		if fileExt(objPath) != extO {
			objPath = output + extO
		}
		err = compileLLWithClang(clangPath, llPath, objPath, verbose)
		if err != nil {
			vprintf("compile ll failed: " + err.Error() + "\n")
			vprintf("You can try manually: " + clangPath + " -c -o " + objPath + " " + llPath + "\n")
			exit(1)
		}
		err = buildRuntimeObjs(clangPath, rtEntries, verbose)
		if err != nil {
			vprintf("building runtime obj failed: " + err.Error() + "\n")
			exit(1)
		}
		err = linkExecutable(clangPath, output, objPath, rtEntries, verbose)
		if err != nil {
			vprintf("link failed: " + err.Error() + "\n")
			vprintf("You can try manually: " + clangPath + " -o " + output + " " + objPath + " runtime/print.o runtime/map.o runtime/io.o\n")
			exit(1)
		}
		if verbose {
			vprintf("built executable " + output + "\n")
		}
	}
}

func fileExt(path string) string {
	for i := len(path) - 1; i >= 0; i = i - 1 {
		if path[i] == 46 {
			return path[i:]
		}
		if path[i] == 47 || path[i] == 92 {
			break
		}
	}
	return ""
}

func deriveLLPath(output string) string {
	ext := fileExt(output)
	base := output
	if ext != "" {
		base = output[:len(output)-len(ext)]
	}
	extLL := ".ll"
	return base + extLL
}

func compileLLWithClang(cc, llPath, objPath string, verbose bool) error {
	args := []string{"-c", "-o", objPath, llPath}
	if verbose {
		vprintf("running clang to compile ll\n")
	}
	cmd := exec.Command(cc, args...)
	outputRes := combinedOutputWrapper(cmd)
	if outputRes.err != nil {
		return errors.New("clang compile ll: " + outputRes.err.Error() + "\n" + string(outputRes.out))
	}
	return nil
}

func buildRuntimeObjs(cc string, entries []struct {
	src string
	out string
}, verbose bool) error {
	for i := 0; i < len(entries); i = i + 1 {
		entry := entries[i]
		args := []string{"-c", "-o", entry.out, entry.src}
		if verbose {
			vprintf("running clang to build runtime obj\n")
		}
		cmd := exec.Command(cc, args...)
		outputRes := combinedOutputWrapper(cmd)
		if outputRes.err != nil {
			return errors.New("cc runtime: " + outputRes.err.Error() + "\n" + string(outputRes.out))
		}
	}
	return nil
}

func linkExecutable(cc, output, objPath string, rtEntries []struct {
	src string
	out string
}, verbose bool) error {
	args := []string{"-o", output, objPath}
	for i := 0; i < len(rtEntries); i = i + 1 {
		args = append(args, rtEntries[i].out)
	}
	if verbose {
		vprintf("running clang to link\n")
	}
	cmd := exec.Command(cc, args...)
	outputRes := combinedOutputWrapper(cmd)
	if outputRes.err != nil {
		return errors.New("link: " + outputRes.err.Error() + "\n" + string(outputRes.out))
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
