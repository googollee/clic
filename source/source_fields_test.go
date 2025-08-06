package source

import (
	"reflect"

	"github.com/googollee/clic/structtags"
)

var (
	a1, a2, a3 string

	fields = []structtags.Field{
		{Name: []string{"a1"}, Description: "a1", Parser: parserString, Value: reflect.ValueOf(&a1).Elem()},
		{Name: []string{"l1", "a2"}, Description: "a2", Parser: parserString, Value: reflect.ValueOf(&a2).Elem()},
		{Name: []string{"l2", "l3", "a3"}, Description: "a3", Parser: parserString, Value: reflect.ValueOf(&a3).Elem()},
	}
)
