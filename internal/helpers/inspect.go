package helper

import (
	"fmt"
	"reflect"
)

func Inspect(object interface{}) {
	fooType := reflect.TypeOf(object)
	fmt.Println("inspect: %s", fooType)
	for i := 0; i < fooType.NumMethod(); i++ {
		method := fooType.Method(i)
		fmt.Println(method.Name)
	}
}
