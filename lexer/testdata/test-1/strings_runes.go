package main

func stringZoo() string {
	interp := "hello\tworld\n"
	raw := `raw
multiline
string`

	var letter rune = 'A'
	var escaped rune = '\n'
	var unicodePi rune = 'π'
	_ = letter
	_ = escaped
	_ = unicodePi

	naïve := "unicode identifier"
	π := 3.14159
	_ = π

	return interp + raw + naïve
}
