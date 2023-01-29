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

type Cmp[T any] func(l, r T) int

type Chain[T any] []Cmp[T]

func (c Chain[T]) Less(l, r T) bool {
	for _, f := range c {
		switch f(l, r) {
		case -1:
			return true
		case 1:
			return false
		}
	}
	return false
}

func By[C any, F constraints.Ordered](by func(C) F, d Dir) Cmp[C] {
	if d {
		return func(l, r C) int {
			bl, br := by(l), by(r)
			switch {
			case bl < br:
				return -1
			case bl > br:
				return 1
			default:
				return 0
			}
		}
	}
	return func(l, r C) int {
		bl, br := by(l), by(r)
		switch {
		case bl < br:
			return 1
		case bl > br:
			return -1
		default:
			return 0
		}
	}
}

func ByFunc[C, F any](by func(C) F, cmp Cmp[F], d Dir) Cmp[C] {
	if d {
		return func(l, r C) int { return cmp(by(l), by(r)) }
	}
	return func(l, r C) int { return cmp(by(r), by(l)) }
}

func ByBytes[C any](by func(C) []byte, d Dir) Cmp[C] {
	if d {
		return func(l, r C) int { return bytes.Compare(by(l), by(r)) }
	}
	return func(l, r C) int { return bytes.Compare(by(r), by(l)) }
}
