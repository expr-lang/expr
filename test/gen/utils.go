package main

import (
	"math/rand"
)

func maybe() bool {
	return rand.Intn(2) == 0
}

type list[T any] []element[T]

type element[T any] struct {
	value  T
	weight int
}

func oneOf[T any](cases []element[T]) T {
	total := 0
	for _, c := range cases {
		total += c.weight
	}
	r := rand.Intn(total)
	for _, c := range cases {
		if r < c.weight {
			return c.value
		}
		r -= c.weight
	}
	return cases[0].value
}

func random[T any](array []T) T {
	return array[rand.Intn(len(array))]
}
