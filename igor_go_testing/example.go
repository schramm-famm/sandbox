package main

import (
	"time"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type B struct {
	length int
}

func NewBig() *B {
	time.Sleep(1000 * time.Millisecond)
	return &B{7}
}

func (b *B) Len() int {
	return b.length
}

func main() {}
