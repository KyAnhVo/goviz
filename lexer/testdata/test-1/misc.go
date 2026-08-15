package main

func miscDemo() {
	a, b := 1, 2
	a, b = b, a
	_, _ = a, b

	x := 10; y := 20; _ = x + y

	/* block comment */ z := 1 + /* inline */ 2
	_ = z
}
