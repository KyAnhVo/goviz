package main

import "fmt"

type Animal struct {
	Name string `json:"name"`
}

func (a Animal) Speak() string {
	return a.Name + " makes a sound"
}

type Dog struct {
	Animal
	Breed string
}

func structDemo() string {
	d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Labrador"}

	point := struct {
		X, Y int
	}{X: 1, Y: 2}

	adder := func(a, b int) int {
		return a + b
	}

	sum := func() int {
		return adder(point.X, point.Y)
	}()

	return fmt.Sprintf("%s (%d,%d) sum=%d", d.Speak(), point.X, point.Y, sum)
}
