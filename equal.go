package sme

import "sync"

// IsEqual compare two sync.Map is equal or not.
func IsEqual(mu *sync.Mutex, a *sync.Map, b *sync.Map) bool {
	if mu == nil {
		mu = &sync.Mutex{}
	}
	mu.Lock()
	defer mu.Unlock()

	isSame := true
	a.Range(func(k, v interface{}) bool {
		bV, ok := b.Load(k)
		if bV != v || !ok {
			isSame = false
			return false
		}
		return true
	})
	b.Range(func(k, v interface{}) bool {
		aV, ok := a.Load(k)
		if aV != v || !ok {
			isSame = false
			return false
		}
		return true
	})
	return isSame
}
