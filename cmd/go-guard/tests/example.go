package main

import "io"

// Str type.
type Str struct{}

type arg struct{}

type other interface{}

// NewStr type
func NewStr(r io.Reader, o other, a *arg, m map[string]string, x string) *Str {
	s := Str{}
	return &s
}
