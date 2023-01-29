package mero

import (
	"bytes"

	"golang.org/x/exp/constraints"
)

const (
	Asc  Dir = true
	Desc Dir = false
)

type Dir bool

type Less[T any] func(l, r T) bool

type Chain[T any] []Less[T]

func (c Chain[T]) Less(l, r T) bool {
	for _, f := range c {
		switch {
		case f(l, r):
			return true
		case f(r, l):
			return false
		}
	}
	return false
}

func By[C any, F constraints.Ordered](by func(C) F, d Dir) Less[C] {
	if d {
		return func(l, r C) bool { return by(l) < by(r) }
	}
	return func(l, r C) bool { return by(r) < by(l) }
}

func ByFunc[C, F any](by func(C) F, less Less[F], d Dir) Less[C] {
	if d {
		return func(l, r C) bool { return less(by(l), by(r)) }
	}
	return func(l, r C) bool { return less(by(r), by(l)) }
}

func ByBytes[C any](by func(C) []byte, d Dir) Less[C] {
	if d {
		return func(l, r C) bool { return bytes.Compare(by(l), by(r)) < 0 }
	}
	return func(l, r C) bool { return bytes.Compare(by(r), by(l)) < 0 }
}
