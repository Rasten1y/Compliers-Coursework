package main

// Простые обёртки над рантаймом.
func print(s string) {
	b := []byte(s)
	if len(b) == 0 {
		var z byte
		gominic_print(&z, 0)
		return
	}
	gominic_print(&b[0], int64(len(b)))
}

func printInt(v int64) {
	gominic_printInt(v)
}

func println() {
	gominic_println()
}

var g1 int = 1
var g2 = 2
var gstr string = "global"
var garr = [2]int{1, 2}
var gstruct = struct {
	x int
	y string
}{x: 3, y: "q"}

func add(a int, b int) int {
	return a + b
}

func main() {
	a := 1
	var b int = 2
	c := add(a, b)

	if c < 10 {
		c = c * 2
	} else {
		c = c - 1
	}

	for i := 0; i < 3; i = i + 1 {
		c = c + i
	}

	printInt(int64(c))
	println()

	for c > 0 {
		c = c - 1
		if c == 0 {
			break
		}
		if c < 0 {
			continue
		}
	}

	msg := "hi " + "there"
	print(msg)
	println()

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

	st := struct {
		a int
		b string
	}{a: c, b: msg}
	_ = st
	print("after struct\n")

	f := float64(1) + 2.5
	if f > 0 {
		print("float ok\n")
	}
	print("after float\n")

	printInt(int64(c))
	println()
	print("end\n")
}
