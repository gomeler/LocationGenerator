package main

import "fmt"

type Animals interface {
	animals []Animal
}

type Animal interface {
	Speak() string
}

type Howlers struct{}

func (h Howlers) Speak() string {
	return "HOWWWWWWLLL"
}

type Dog struct{}

func (d Dog) Speak() string {
	return "Woof!"
}

type Wolf struct {
	Howlers
}

type Beagle struct {
	Howlers
}

type Cat struct{}

func (c Cat) Speak() string {
	return "Meow"
}

func main() {
	var a := Wolf{}
	var b := Wolf{}
	var farm = Animals{a, b}
	
}
