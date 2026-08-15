package main

func numberZoo() float64 {
	dec := 1_000_000
	hex := 0x1A_F
	oct := 0o17
	bin := 0b1010_1010
	flt := 1.5e3
	small := .25
	imag := 2 + 3i
	_ = imag

	return float64(dec) + float64(hex) + float64(oct) + float64(bin) + flt + small
}
