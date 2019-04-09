package gotest

import (
	"fmt"
	"reflect"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/mock"
)

type Test struct {
	mock.Mock
}

// Stub a func to this Test.
//
//     t.Stub("MyMethod", Func)
func (t *Test) StubFunc(fnIn, fnOut interface{}) {
	monkey.Patch(fnIn, fnOut)
	return
}

// StubInstanceMethod a func to this Test.
//
//     var d *net.Dialer
//     t.StubInstFunc(d, "Dial", Func)
func (t *Test) StubInstFunc(target interface{}, methodName string, replacement interface{}) {
	monkey.PatchInstanceMethod(reflect.TypeOf(target), methodName, replacement)
	return
}

// Mock a func to this Test.
//
//     t.Mock("MyMethod", Func)
func (t *Test) MockFunc(methodName string, fn interface{}) {
	if v := reflect.ValueOf(fn); v.Kind() != reflect.Func {
		panic(fmt.Sprintf("must be a Func in expectations. fn Type is \"%T\")", fn))
	}
	ft := reflect.TypeOf(fn)
	mfn := reflect.MakeFunc(ft, func(args []reflect.Value) (results []reflect.Value) {
		vargs := []interface{}{}
		for i := range args {
			vargs = append(vargs, args[i].Interface())
		}
		ret := t.MethodCalled(methodName, vargs...)
		for i := 0; i < reflect.TypeOf(fn).NumOut(); i++ {
			results = append(results, reflect.ValueOf(ret.Get(i)))
		}
		return

	})
	monkey.Patch(fn, mfn.Interface())
	return
}

// MockInstFunc a func to this Test.
//
//     var d *net.Dialer
//     t.MockInstFunc("MyMethod", d, "Dial")
func (t *Test) MockInstFunc(methodName string, target interface{}) {
	tf := reflect.TypeOf(target)
	mtd, ok := tf.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("must be a Func in expectations. fn Type is \"%T\")", target))
	}
	ft := mtd.Type
	if ft.Kind() != reflect.Func {
		panic(fmt.Sprintf("must be a Func in expectations. fn Type is \"%T\")", ft))
	}
	mfn := reflect.MakeFunc(mtd.Type, func(args []reflect.Value) (results []reflect.Value) {
		vargs := []interface{}{}
		for i := range args {
			vargs = append(vargs, args[i].Interface())
		}
		ret := t.MethodCalled(methodName, vargs...)
		for i := 0; i < ft.NumOut(); i++ {
			results = append(results, reflect.ValueOf(ret.Get(i)))
		}
		return

	})
	monkey.PatchInstanceMethod(tf, methodName, mfn.Interface())
	return
}

// Call at the end of this Test.
//
//     t.Close()
func (t *Test) Close() {
	monkey.UnpatchAll()
	return
}
