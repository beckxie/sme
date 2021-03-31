package sme_test

import (
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/beckxie/sme"
)

var mu sync.Mutex
var syncMapA, syncMapB *sync.Map

func initSyncMap() {
	const maxI = 10000

	syncMapA = &sync.Map{}
	syncMapB = &sync.Map{}

	for i := 0; i < maxI; i++ {
		syncMapA.Store(strconv.Itoa(i), i)
		syncMapB.Store(strconv.Itoa(i), i)
	}
}

func isEqualByTwoMap(mu *sync.Mutex, a *sync.Map, b *sync.Map) bool {
	mu.Lock()
	defer mu.Unlock()

	aMap := make(map[interface{}]interface{})
	bMap := make(map[interface{}]interface{})

	// fmt.Println("a:")
	a.Range(func(k, v interface{}) bool {
		// fmt.Printf("k:%#v , v:%#v\n", k, v)
		aMap[k] = v
		return true
	})
	// fmt.Println("---")
	// fmt.Println("b:")
	b.Range(func(k, v interface{}) bool {
		// fmt.Printf("k:%#v , v:%#v\n", k, v)
		bMap[k] = v
		return true
	})

	return reflect.DeepEqual(aMap, bMap)
}

func BenchmarkIsEqualByPointerRange(b *testing.B) {
	initSyncMap()
	for i := 0; i < b.N; i++ {
		sme.IsEqual(&mu, syncMapA, syncMapB)
	}
}

func BenchmarkIsEqualByTwoMap(b *testing.B) {
	initSyncMap()
	for i := 0; i < b.N; i++ {
		isEqualByTwoMap(&mu, syncMapA, syncMapB)
	}
}
