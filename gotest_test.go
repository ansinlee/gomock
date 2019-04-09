package gotest

import (
	"fmt"
	"testing"
)

type T struct {
	A int
}

func (t *T) Dosomething(a int) int {
	return a + t.A
}

func Dosomething(a int) int {
	t := &T{1}
	v := t.Dosomething(a)
	return v
}

func TestStubFunc(t *testing.T) {
	tt := new(Test)
	defer tt.Close()

	tt.StubFunc(Dosomething, func(a int) int {
		fmt.Println("stub Dosomething")
		return a + 100
	})

	if Dosomething(1) != 101 {
		t.Fatal("stub Dosomething failed")
	}
}

func TestStubInstFunc(t *testing.T) {
	tt := new(Test)
	defer tt.Close()

	st := &T{A: 1}

	tt.StubInstFunc(st, "Dosomething", func(_ *T, a int) int {
		fmt.Println("stub T.Dosomething")
		return st.A + a + 1
	})

	if st.Dosomething(1) != 3 {
		t.Fatal("stub T.Dosomething failed")
	}
}

func TestMockFunc(t *testing.T) {
	tt := new(Test)
	defer tt.Close()

	// mock
	tt.MockFunc("Dosomething", Dosomething)

	//collaborator
	tt.On("Dosomething", 1).Return(2)

	//test
	Dosomething(1)

	//assert
	tt.AssertExpectations(t)
}

func TestMockInstFunc(t *testing.T) {
	tt := new(Test)
	defer tt.Close()

	// mock
	tt.MockInstFunc("Dosomething", new(T))

	//collaborator
	tt.On("Dosomething", &T{1}, 1).Return(2)

	//test
	Dosomething(1)

	//assert
	tt.AssertExpectations(t)
}

func TestPatchValue(t *testing.T) {
	value := int(1)

	tt := new(Test)

	tt.PatchValue(&value, 2)

	if value != 2 {
		t.Fatal("patch value failed")
	}
	tt.Close()

	if value != 1 {
		t.Fatal("recover patch value failed")
	}
}
