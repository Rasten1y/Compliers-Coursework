package frontend

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func parseSnippet(t *testing.T, src string) *ast.File {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "snippet.go", src, 0)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	return file
}

func TestCheckUnsupportedAllowsMinimalSubset(t *testing.T) {
	src := `package main

var g int

func main() {
	x := 1
	y := 2
	if x < y {
		x++
	} else {
		for i := 0; i < 3; i++ {
			y = y + i
			if y > 5 {
				break
			}
			continue
		}
	}
	_ = x
	_ = y
}`
	file := parseSnippet(t, src)
	if err := CheckUnsupported(file); err != nil {
		t.Fatalf("unexpected unsupported construct: %v", err)
	}
}

func TestCheckUnsupportedRejectsUnsupportedNodes(t *testing.T) {
	cases := []struct {
		name string
		src  string
	}{
		{"go", `package main; func main(){ go foo() }`},
		{"defer", `package main; func main(){ defer foo() }`},
		{"select", `package main; func main(){ select{} }`},
		{"switch", `package main; func main(){ switch 1 {} }`},
		{"chan", `package main; func main(){ var ch chan int; ch <- 1 }`},
		{"func lit", `package main; func main(){ f := func(){}; _ = f }`},
		{"goto", `package main; func main(){ goto L; L: return }`},
		{"labeled", `package main; func main(){ L: _ = 1 }`},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			file := parseSnippet(t, tc.src)
			if err := CheckUnsupported(file); err == nil {
				t.Fatalf("expected unsupported construct, got none")
			}
		})
	}
}
