package utils

import (
	"fmt"
	"reflect"
)

func Required[T any](value T, name string) T {
	if reflect.ValueOf(value).IsZero() {
		panic(fmt.Sprintf("%s is required", name))
	}
	return value
}
