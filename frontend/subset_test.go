package frontend

import (
	"os"
	"path/filepath"
	"testing"
)

// Test that a sample program using our supported subset parses, type-checks and lowers to IR.
func TestSubsetCoverage(t *testing.T) {
	src := `package main

var g1 int = 1
var g2 = 2
var gstr string = "global"
var garr = [2]int{1, 2}
var gstruct = struct { x int; y string }{ x: 3, y: "q" }

func add(a int, b int) int { return a + b }

func main() {
	var a int = 1
	b := 2
	c := add(a, b)
	if c < 10 {
		c = c * 2
	} else {
		c = c - 1
	}
	for i := 0; i < 3; i = i + 1 {
		c = c + i
	}
	for c > 0 {
		c = c - 1
		if c == 0 { break }
		if c < 0 { continue }
	}
	s := "hi " + "there"
	_ = s
	arr := [2]int{1, 2}
	_ = arr[1]
	sl := []int{}
	sl = append(sl, 5)
	_ = len(sl)
	m := make(map[int]int)
	m[1] = c
	_, ok := m[1]
	if ok {
		c = c + m[1]
	}
	st := struct { a int; b string }{ a: c, b: s }
	_ = st
	f := float64(1) + 2.5
	_ = f
}
`
	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	SetSkipSubsetCheck(false)
	progs, err := ParseAndCheckAll([]string{path})
	if err != nil {
		t.Fatalf("ParseAndCheckAll failed: %v", err)
	}
	if _, err := BuildModule(progs); err != nil {
		t.Fatalf("BuildModule failed: %v", err)
	}
}
