package main

import "fmt"

type A struct {
}

func (a A) test() {
	fmt.Println("test")
}

func (a *A) test2() {
	fmt.Println("test2")
}

func main() {
	a := &A{}
	(*A).test(a)
	(*A).test2(a)
}
