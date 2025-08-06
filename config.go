package clic

import (
	"context"
	"fmt"
	"reflect"
)

type configAdapter interface {
	Callback(context.Context) error
	Value() reflect.Value
}

type configAny struct {
	value    reflect.Value
	callback reflect.Value
}

func newConfigValue(value any) *configAny {
	if value == nil {
		panic("register with nil value")
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("register with invalid type %T, must be `*Type`", value))
	}

	return &configAny{
		value: v,
	}
}

func newConfigCallback(callback any) *configAny {
	if callback == nil {
		panic("register with nil callback")
	}

	f := reflect.ValueOf(callback)
	if f.Kind() != reflect.Func {
		panic(fmt.Sprintf("register with invalid callback %T, must be `func(context.Context, *Type) error`", callback))
	}

	ft := f.Type()
	if ft.NumIn() != 2 || ft.NumOut() != 1 {
		panic(fmt.Sprintf("register with invalid callback %T, must be `func(context.Context, *Type) error`", callback))
	}

	if ft.In(0) != reflect.TypeFor[context.Context]() {
		panic(fmt.Sprintf("register with invalid callback %T, must be `func(context.Context, *Type) error`", callback))
	}

	valType := ft.In(1)
	if valType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("register with invalid callback %T, must be `func(context.Context, *Type) error`", callback))
	}

	if ft.Out(0) != reflect.TypeFor[error]() {
		panic(fmt.Sprintf("register with invalid callback %T, must be `func(context.Context, *Type) error`", callback))
	}

	return &configAny{
		value:    reflect.New(valType.Elem()),
		callback: f,
	}
}

func (c *configAny) Callback(ctx context.Context) error {
	if !c.callback.IsValid() {
		return nil
	}

	in := []reflect.Value{reflect.ValueOf(ctx), c.value}
	out := c.callback.Call(in)

	if out[0].IsNil() {
		return nil
	}

	return out[0].Interface().(error)
}

func (c *configAny) Get() any {
	return c.value.Interface()
}

func (c *configAny) Value() reflect.Value {
	return c.value
}
