package main

import (
	"math"
	"fmt"
)

type person struct {
	name string
}

type secretAgent struct {
	person
	agent bool
}

func (p person) speak() string {
	return fmt.Sprintln("Hello, I'm a person and my name is", p.name)
}

func (sa secretAgent) speak() string {
	return fmt.Sprintf("Hello, i'm an agent. My name is %s and I have a licence to kill: %v\n",
		sa.name, sa.agent)
}

type allPersons interface {
	speak() string
}


type square struct {
	side float64
}

func (s square) area() float64 {
	return math.Pow(s.side, 2)
}

type circle struct {
	radius float64
}

func (c circle) area() float64 {
	return math.Pi * math.Pow(c.radius, 2)
}

type shape interface {
	area() float64
}

func info(s shape) {
	fmt.Println("from the shape interface", s.area())
}

func vomit(a allPersons) {
	switch v := a.(type) {
	case person:
		fmt.Println(v.name)
	case secretAgent:
		fmt.Println(v.name)
	default:
		fmt.Println("unknown")

	}
}

func main() {
	s1 := square{3}
	c1 := circle{2}
	fmt.Println(s1.area())
	fmt.Println(c1.area())

	info(s1)
	info(c1)

	p1 := person{
		name:"Ross",
	}
	sa1 := secretAgent{
		person{"James"},
		true,
	}

	fmt.Println(p1.name)
	p1.speak()

	fmt.Println(sa1.name)
	fmt.Println(sa1.speak())
	fmt.Println(sa1.person.speak())

	vomit(p1)
	vomit(sa1)
}