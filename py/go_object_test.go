package py

import (
	"fmt"
	"testing"
)

type MyStruct struct {
	Doc    string `gpython:"__doc__"`
	Field1 int    `gpython:"field1"`
	Field2 string
}

func (s MyStruct) Method1(a int, b string) {
	fmt.Println("Method1 called with arguments:", a, b)
}

func (s MyStruct) Method_plus(a, b int) int {
	fmt.Println("Method_plus called")
	return a + b
}

func (s MyStruct) Method_method() {
	fmt.Println("Method_method called")
}

func (s MyStruct) Property_value() string {
	fmt.Println("Property_value called")
	return "TestProperty"
}

func TestNewGoObjectFrom(t *testing.T) {
	s := &MyStruct{
		Doc:    "this is a doc",
		Field1: 666,
		Field2: "test",
	}
	_ = s
	gs := NewGoObject(nil)
	v, err := gs.GetField("__doc__")
	if err != nil {
		panic(err)
	}
	fmt.Println(v)

	v, err = gs.GetField("field1")
	if err != nil {
		panic(err)
	}
	fmt.Println(v)

	err = gs.SetField("field1", 999)
	if err != nil {
		panic(err)
	}
	v, err = gs.GetField("field1")
	if err != nil {
		panic(err)
	}
	fmt.Println(v)

	res, err := gs.CallMethod("method")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	res, err = gs.CallMethod("plus", 100, 666)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	res, err = gs.CallMethod("value")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
