package clic

import (
	"context"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	type cfg struct {
		Field1 int
		Field2 string
	}
	var callbacked bool

	tests := []struct {
		callback       func(context.Context, *cfg) error
		wantCallbacked bool
	}{
		{func(context.Context, *cfg) error {
			callbacked = true
			return nil
		}, true},
		{nil, false},
	}

	for _, tc := range tests {
		var c config[cfg]
		c.callback = tc.callback

		if got, want := c.Value().Type(), reflect.TypeFor[cfg](); got != want {
			t.Errorf("c.Value().Type() = %v, want: %v", got, want)
		}

		callbacked = false
		if err := c.Callback(context.TODO()); err != nil {
			t.Errorf("c.Callback() = %v, want: no error", err)
		}
		if got, want := callbacked, tc.wantCallbacked; got != want {
			t.Errorf("callbacked = %v, want: %v", got, want)
		}

		if c.Get() == nil {
			t.Errorf("c.Get() == nil, which should not")
		}
	}
}
