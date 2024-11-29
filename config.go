package clic

import (
	"context"
	"reflect"
)

type config[Config any] struct {
	value    Config
	callback func(context.Context, *Config) error
}

func (c *config[Config]) Callback(ctx context.Context) error {
	if c.callback == nil {
		return nil
	}

	return c.callback(ctx, &c.value)
}

func (c *config[Config]) Get() *Config {
	return &c.value
}

func (c *config[Config]) Value() reflect.Value {
	return reflect.ValueOf(&c.value).Elem()
}

type configHandler interface {
	Callback(context.Context) error
	Value() reflect.Value
}

var handlers = map[string]configHandler{}
