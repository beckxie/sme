package main

import (
	"fmt"
	"sync"

	"github.com/beckxie/sme"
)

func main() {
	a := &sync.Map{}
	a.Store(1, nil)
	a.Store("2", 2)
	a.Store("3", "3")

	b := &sync.Map{}
	b.Store(1, nil)
	b.Store("2", 2)
	b.Store("3", "3")

	c := &sync.Map{}
	c.Store(1, nil)
	c.Store("2", 2)
	c.Store("3", 3)

	var mu sync.Mutex
	fmt.Printf("Simple compare sync.Map-a and sync.Map-b result: %t\n", sme.IsEqual(&mu, a, b))
	fmt.Printf("Simple compare sync.Map-a and sync.Map-c result: %t\n", sme.IsEqual(&mu, a, c))
}
