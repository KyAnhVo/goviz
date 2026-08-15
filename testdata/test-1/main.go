package main

import "fmt"

func main() {
	fmt.Println("numbers:", numberZoo())
	fmt.Println("strings:", stringZoo())
	fmt.Println("generic max:", Max(3, 7))
	fmt.Println("control flow:", controlFlow())
	concurrencyDemo()
	fmt.Println("structs:", structDemo())
	miscDemo()
}
