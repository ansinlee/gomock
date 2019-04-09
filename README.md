GoTest 
================================

Go code (golang) set of packages that provide many tools for testifying that your code will behave as you intend.

Features include:

  * [Stub] stub both package and struct method
  * [Mock] mock package and stuct method without interface

Get started:

  * Install testify with [one line of code](#installation), or [update it with another](#staying-up-to-date)
  * For an introduction to writing test code in Go, see http://golang.org/doc/code.html#Testing
  * Check out the API Documentation http://godoc.org/github.com/ansinlee/gotest
  * A little about [Test-Driven Development (TDD)](http://en.wikipedia.org/wiki/Test-driven_development)



[`gotest`](http://godoc.org/github.com/ansinlee/gotest "API documentation") package
-------------------------------------------------------------------------------------------

The `assert` package provides some helpful methods that allow you to write better test code in Go.

  * Prints friendly, easy to read failure descriptions
  * Allows for very readable code
  * Optionally annotate each assertion with a message

See it in action:
```go
// run command: go test -v -gcflags=-l

package yours

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
```

For more information on how to write mock code, check out the [API documentation for the `gotest` package](http://godoc.org/github.com/ansinlee/gotest).
------

Installation
============

To install Testify, use `go get`:

    go get github.com/ansinlee/gotest

This will then make the following packages available to you:

    github.com/stretchr/testify/mock
    github.com/bouk/monkey

------

Staying up to date
==================

To update Testify to the latest version, use `go get -u github.com/ansinlee/gotest`.

------

Supported go versions
==================

We support the three major Go versions, which are 1.9, 1.10, and 1.11 at the moment.

------

Notes
==================

1. Monkey sometimes fails to patch a function if inlining is enabled. Try running your tests with inlining disabled, for example: `go test -gcflags=-l`. The same command line argument can also be used for build.
2. Monkey won't work on some security-oriented operating system that don't allow memory pages to be both write and execute at the same time. With the current approach there's not really a reliable fix for this.
3. Monkey is not threadsafe. Or any kind of safe.
4. I've tested monkey on OSX 10.10.2 and Ubuntu 14.04. It should work on any unix-based x86 or x86-64 system.

------

Contributing
============

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue.  Extra credit for those using Testify to write the test code that demonstrates it.

------

License
=======

This project is licensed under the terms of the MIT license.
