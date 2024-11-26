package structtags

import (
	"fmt"
	"reflect"
	"testing"
)

type TestWithChan struct {
	Chan chan int
}

type TestLayerWithChan struct {
	Inner TestWithChan `clic:"inner"`
}

func TestStructParseInvalid(t *testing.T) {
	var chanValue TestWithChan
	var layerValue TestLayerWithChan

	tests := []struct {
		value any
	}{
		{&chanValue},
		{&layerValue},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%T", tc.value), func(t *testing.T) {
			for _, value := range []reflect.Value{reflect.ValueOf(tc.value).Elem(), reflect.ValueOf(tc.value)} {
				_, err := ParseStruct(value, []string{"test"})
				if err == nil {
					t.Errorf("ParseStruct(%v, ['test']) should returns an error, but not", value)
				}
			}
		})
	}
}
