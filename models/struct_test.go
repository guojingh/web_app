package models

import (
	"fmt"
	"testing"
	"unsafe"
)

type s1 struct {
	a int8
	b string
	c int8
}

type s2 struct {
	a int8
	b int8
	c string
}

func TestStruct(t *testing.T) {
	v1 := s1{
		a: 1,
		b: "guojinghu",
		c: 2,
	}
	v2 := s2{
		a: 1,
		b: 2,
		c: "guojinghu",
	}

	fmt.Println(unsafe.Sizeof(v1), unsafe.Sizeof(v2))
}
