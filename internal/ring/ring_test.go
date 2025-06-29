package ring

import (
	"fmt"
	"testing"
)

func TestRing(t *testing.T) {
	type op = ringOp[int]
	testRing(t, New[int](3),
		// noops on empty ring
		op{cap: 0, opType: opRst, value: 0, items: []int{}},
		op{cap: 0, opType: opDeq, value: 0, items: []int{}},

		// basic
		op{cap: 3, opType: opEnq, value: 1, items: []int{1}},
		op{cap: 3, opType: opDeq, value: 1, items: []int{}},

		// wrapping
		op{cap: 3, opType: opEnq, value: 2, items: []int{2}},
		op{cap: 3, opType: opEnq, value: 3, items: []int{2, 3}},
		op{cap: 3, opType: opEnq, value: 4, items: []int{2, 3, 4}},
		op{cap: 3, opType: opDeq, value: 2, items: []int{3, 4}},
		op{cap: 3, opType: opDeq, value: 3, items: []int{4}},
		op{cap: 3, opType: opDeq, value: 4, items: []int{}},

		// resetting
		op{cap: 3, opType: opEnq, value: 2, items: []int{2}},
		op{cap: 3, opType: opRst, value: 0, items: []int{}},
		op{cap: 3, opType: opDeq, value: 0, items: []int{}},

		// growing without wrapping
		op{cap: 3, opType: opEnq, value: 5, items: []int{5}},
		op{cap: 3, opType: opEnq, value: 6, items: []int{5, 6}},
		op{cap: 3, opType: opEnq, value: 7, items: []int{5, 6, 7}},
		op{cap: 6, opType: opEnq, value: 8, items: []int{5, 6, 7, 8}},
		op{cap: 6, opType: opRst, value: 0, items: []int{}},
		op{cap: 6, opType: opDeq, value: 0, items: []int{}},

		// growing and wrapping
		op{cap: 6, opType: opEnq, value: 9, items: []int{9}},
		op{cap: 6, opType: opEnq, value: 10, items: []int{9, 10}},
		op{cap: 6, opType: opEnq, value: 11, items: []int{9, 10, 11}},
		op{cap: 6, opType: opEnq, value: 12, items: []int{9, 10, 11, 12}},
		op{cap: 6, opType: opEnq, value: 13, items: []int{9, 10, 11, 12, 13}},
		op{cap: 6, opType: opEnq, value: 14, items: []int{9, 10, 11, 12, 13, 14}},
		op{cap: 6, opType: opDeq, value: 9, items: []int{10, 11, 12, 13, 14}},
		op{cap: 6, opType: opDeq, value: 10, items: []int{11, 12, 13, 14}},
		op{cap: 6, opType: opEnq, value: 15, items: []int{11, 12, 13, 14, 15}},
		op{cap: 6, opType: opEnq, value: 16, items: []int{11, 12, 13, 14, 15, 16}},
		op{cap: 9, opType: opEnq, value: 17, items: []int{11, 12, 13, 14, 15, 16, 17}}, // grows wrapped
		op{cap: 9, opType: opDeq, value: 11, items: []int{12, 13, 14, 15, 16, 17}},
		op{cap: 9, opType: opDeq, value: 12, items: []int{13, 14, 15, 16, 17}},
		op{cap: 9, opType: opDeq, value: 13, items: []int{14, 15, 16, 17}},
		op{cap: 9, opType: opDeq, value: 14, items: []int{15, 16, 17}},
		op{cap: 9, opType: opDeq, value: 15, items: []int{16, 17}},
		op{cap: 9, opType: opDeq, value: 16, items: []int{17}},
		op{cap: 9, opType: opDeq, value: 17, items: []int{}},
		op{cap: 9, opType: opDeq, value: 0, items: []int{}},
	)

	t.Run("should panic on invalid chunkSize", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("should have panicked")
			}
		}()
		New[int](0)
	})
}

const (
	opEnq = iota // enqueue an item
	opDeq        // dequeue an item and an item was available
	opRst        // reset
)

type ringOp[T comparable] struct {
	cap    int // expected values
	opType int // opEnq or opDeq
	value  T   // value to enqueue or value expected for dequeue; ignored for opRst
	items  []T // items left
}

func testRing[T comparable](t *testing.T, r *Ring[T], ops ...ringOp[T]) {
	for i, op := range ops {
		testOK := t.Run(fmt.Sprintf("opIndex=%v", i), func(t *testing.T) {
			testRingOp(t, r, op)
		})
		if !testOK {
			return
		}
	}
}

func testRingOp[T comparable](t *testing.T, r *Ring[T], op ringOp[T]) {
	var zero T
	switch op.opType {
	case opEnq:
		r.Enqueue(op.value)
	case opDeq:
		shouldSucceed := r.Len() > 0
		v, ok := r.Dequeue()
		switch {
		case ok != shouldSucceed:
			t.Fatalf("should have succeeded: %v", shouldSucceed)
		case ok && v != op.value:
			t.Fatalf("expected value: %v; got: %v", op.value, v)
		case !ok && v != zero:
			t.Fatalf("expected zero value; got: %v", v)
		}
	case opRst:
		r.Reset()
	}
	if c := r.Cap(); c != op.cap {
		t.Fatalf("expected cap: %v; got: %v", op.cap, c)
	}
	if l := r.Len(); l != len(op.items) {
		t.Errorf("expected Len(): %v; got: %v", len(op.items), l)
	}
	var got []T
	for i := 0; ; i++ {
		v, ok := r.Nth(i)
		if !ok {
			break
		}
		got = append(got, v)
	}
	if l := len(got); l != len(op.items) {
		t.Errorf("expected items: %v\ngot items: %v", op.items, got)
	}
	for i := range op.items {
		if op.items[i] != got[i] {
			t.Fatalf("expected items: %v\ngot items: %v", op.items, got)
		}
	}
	if v, ok := r.Nth(len(op.items)); ok || v != zero {
		t.Fatalf("expected no more items, got: v=%v; ok=%v", v, ok)
	}
}
