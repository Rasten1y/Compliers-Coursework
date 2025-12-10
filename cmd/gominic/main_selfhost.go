//go:build selfhost
// +build selfhost

package main

import (
	"gominic/backend"
	"gominic/frontend"
)

// Self-hosted entry point that avoids os/exec and file I/O from the Go stdlib.
// It relies solely on runtime helpers (gominic_argc/argv/write_file/print).

func main() {
	argc := gominic_argc()
	if argc < 2 {
		printString("usage: gominic_selfhost [-o out.ll] files...\n")
		return
	}

	output := "a.ll"
	var files []string
	for i := int64(1); i < argc; i++ {
		arg := gominic_argv(i)
		if arg == "-o" && i+1 < argc {
			output = gominic_argv(i + 1)
			i++
			continue
		}
		files = append(files, arg)
	}

	progs, err := frontend.ParseAndCheckAll(files)
	if err != nil {
		printString("frontend failed: " + err.Error() + "\n")
		return
	}

	mod, err := frontend.BuildModule(progs)
	if err != nil {
		printString("IR lowering failed: " + err.Error() + "\n")
		return
	}

	irText := backend.EmitModule(mod)
	if !writeFile(output, irText) {
		printString("write ll failed\n")
		return
	}
}

// writeFile writes data to path using runtime helper.
func writeFile(path string, data string) bool {
	pb := []byte(path)
	db := []byte(data)
	var pPtr *byte
	if len(pb) == 0 {
		var z byte
		pPtr = &z
	} else {
		pPtr = &pb[0]
	}
	var dPtr *byte
	if len(db) == 0 {
		var z byte
		dPtr = &z
	} else {
		dPtr = &db[0]
	}
	return gominic_write_file(pPtr, int64(len(pb)), dPtr, int64(len(db)))
}

// printString emits to stdout via runtime print.
func printString(s string) {
	b := []byte(s)
	if len(b) == 0 {
		var z byte
		gominic_print(&z, 0)
		return
	}
	gominic_print(&b[0], int64(len(b)))
}
