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

	var mu sync.Mutex
	fmt.Printf("Simple compare two sync.Map result: %t\n", sme.IsEqual(&mu, a, b))
}
