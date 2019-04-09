package gotest

import (
	"fmt"
	"reflect"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock // mock API Document: https://godoc.org/github.com/stretchr/testify/mock
	patch     map[reflect.Value]reflect.Value
}

// Patch a value to this Mock.
//
//     t.PatchValue(&orm.Debug, true)
func (t *Mock) PatchValue(target, replace interface{}) {
	tv := reflect.Indirect(reflect.ValueOf(target))
	rv := reflect.Indirect(reflect.ValueOf(replace))

	if !tv.CanSet() {
		panic("target has to be a prt can set")
	}

	if t.patch == nil {
		t.patch = make(map[reflect.Value]reflect.Value)
	}

	// only save the oldest
	if _, ok := t.patch[tv]; !ok {
		// copy and save old value
		old := reflect.New(tv.Type()).Elem()
		old.Set(tv)
		t.patch[tv] = old
	}

	tv.Set(reflect.Indirect(rv))
}

// Stub a func to this Mock.
//
//     t.StubFunc("MyMethod", Func)
func (t *Mock) StubFunc(fnIn, fnOut interface{}) {
	monkey.Patch(fnIn, fnOut)
	return
}

// StubInstanceMethod a func to this Mock.
//
//     var d *net.Dialer
//     t.StubInstFunc(d, "Dial", Func)
func (t *Mock) StubInstFunc(target interface{}, methodName string, replacement interface{}) {
	monkey.PatchInstanceMethod(reflect.TypeOf(target), methodName, replacement)
	return
}

// Mock a func to this Mock.
//
//     t.MockFunc("MyMethod", Func)
func (t *Mock) MockFunc(methodName string, fn interface{}) {
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

// MockInstFunc a func to this Mock.
//
//     var d *net.Dialer
//     t.MockInstFunc("MyMethod", d, "Dial")
func (t *Mock) MockInstFunc(methodName string, target interface{}) {
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

// Call at the end of this Mock.
//
//     t.Close()
func (t *Mock) Close() {
	monkey.UnpatchAll()
	for v, it := range t.patch {
		v.Set(it)
	}
	return
}
