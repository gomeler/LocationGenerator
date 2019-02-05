package main

import "fmt"

/*
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
*/

func main() {
	//var a := Wolf{}
	//var b := Wolf{}
	//var farm = Animals{a, b}
	//Had a little confusion with pointers earlier.
	m := &Mutable{0, 0}
	fmt.Println(m)
	m.StayTheSame()
	fmt.Println(m)
	m.Mutate()
	fmt.Println(m)

}

type Mutable struct {
	a int
	b int
}

func (m Mutable) StayTheSame() {
	m.a = 5
	m.b = 7
}

func (m *Mutable) Mutate() {
	m.a = 5
	m.b = 7
}
