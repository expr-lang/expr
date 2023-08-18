package main

import (
	"math/rand"

	"github.com/antonmedv/expr/ast"
)

func maybe() bool {
	return rand.Intn(2) == 0
}

type fn func(int) ast.Node

type fnWeight struct {
	value  fn
	weight int
}

func weightedRandom(cases []fnWeight) fn {
	totalWeight := 0
	for _, c := range cases {
		totalWeight += c.weight
	}
	r := rand.Intn(totalWeight)
	for _, c := range cases {
		if r < c.weight {
			return c.value
		}
		r -= c.weight
	}
	return cases[0].value
}

type intWeight struct {
	value  int
	weight int
}

func weightedRandomInt(cases []intWeight) int {
	totalWeight := 0
	for _, c := range cases {
		totalWeight += c.weight
	}
	r := rand.Intn(totalWeight)
	for _, c := range cases {
		if r < c.weight {
			return c.value
		}
		r -= c.weight
	}
	return cases[0].value
}
