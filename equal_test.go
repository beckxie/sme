package sme_test

import (
	"context"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/beckxie/sme"
)

func TestIsEqual(t *testing.T) {
	type args struct {
		mu *sync.Mutex
		a  *sync.Map
		b  *sync.Map
	}

	argMutex := &sync.Mutex{}
	sma := &sync.Map{}
	sma.Store(1, nil)
	sma.Store("2", 2)
	sma.Store("3", "3")

	smb := &sync.Map{}
	smb.Store(1, nil)
	smb.Store("2", 2)
	smb.Store("3", "3")

	smc := &sync.Map{}
	smc.Store(1, nil)
	smc.Store("2", 2)
	smc.Store("3", 3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fillSyncMapConcurrency(ctx, cancel)

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case#1",
			args: args{
				mu: nil,
				a:  sma,
				b:  smb,
			},
			want: true,
		},
		{
			name: "case#2",
			args: args{
				mu: argMutex,
				a:  sma,
				b:  smb,
			},
			want: true,
		},
		{
			name: "case#3",
			args: args{
				mu: argMutex,
				a:  sma,
				b:  smc,
			},
			want: false,
		},
		{
			name: "worst-case#1",
			args: args{
				mu: nil,
				a:  nil,
				b:  smb,
			},
			want: false,
		},
		{
			name: "worst-case#2",
			args: args{
				mu: nil,
				a:  sma,
				b:  nil,
			},
			want: false,
		},
		{
			name: "worst-case#3",
			args: args{
				mu: argMutex,
				a:  nil,
				b:  nil,
			},
			want: true,
		},
		{
			name: "concurrency-case#1",
			args: args{
				mu: nil,
				a:  syncMapA,
				b:  syncMapB,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sme.IsEqual(tt.args.mu, tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("IsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func BenchmarkIsEqualByRange(b *testing.B) {
	initSyncMap()
	for i := 0; i < b.N; i++ {
		sme.IsEqual(&mu, syncMapA, syncMapB)
	}
}

func fillSyncMapConcurrency(ctx context.Context, cancel context.CancelFunc) {
	go func(ctx context.Context, cancel context.CancelFunc) {
		for {
			const maxI = 10
			syncMapA = &sync.Map{}
			syncMapB = &sync.Map{}

			for i := 0; i < maxI; i++ {
				select {
				case <-ctx.Done():
					cancel()
					return
				default:
				}
				syncMapA.Store(strconv.Itoa(i), i)
				syncMapB.Store(strconv.Itoa(i), i)
			}
		}
	}(ctx, cancel)
}
func BenchmarkIsEqualByRangeConcurrency(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fillSyncMapConcurrency(ctx, cancel)
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
