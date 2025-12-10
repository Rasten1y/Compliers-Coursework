package main

// Helper wrappers to call runtime-expected signatures.
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

func main() {
	print("hello from gominic\n")
	printInt(42)
	println()

	msg := "example"
	print(msg)
	println()
}
