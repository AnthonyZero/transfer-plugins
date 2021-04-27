package utils

import (
	"reflect"
)

// Contain 判断array/slice/map中是否包含某个元素
func Contain(item interface{}, array interface{}) bool {
	target := reflect.ValueOf(array)
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < target.Len(); i++ {
			if reflect.DeepEqual(item, target.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		if target.MapIndex(reflect.ValueOf(item)).IsValid() {
			return true
		}
	}

	return false
}
