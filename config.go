package clic

import (
	"context"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/googollee/clic/sources"
	"github.com/googollee/clic/structtags"
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

type configAdapter interface {
	Callback(context.Context) error
	Value() reflect.Value
}

var configs = newCLIConfigs()

type cliConfigs struct {
	configs map[string]configAdapter
}

func newCLIConfigs() cliConfigs {
	return cliConfigs{
		configs: make(map[string]configAdapter),
	}
}

func (c *cliConfigs) ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(125)
}

func (c *cliConfigs) Register(name string, adapter configAdapter) error {
	if _, exist := c.configs[name]; exist {
		return fmt.Errorf("already registered a config with name %s", name)
	}
	c.configs[name] = adapter

	return nil
}

func (c *cliConfigs) Prepare(srcs []sources.Source, fset *flag.FlagSet) error {
	var fields []structtags.Field
	for name, handler := range c.configs {
		f, err := structtags.ParseStruct(handler.Value(), []string{name})
		if err != nil {
			return fmt.Errorf("parse config %q error: %w", name, err)
		}
		fields = append(fields, f...)
	}

	for i := 0; i < len(srcs); i++ {
		src := srcs[i]
		if err := src.Prepare(fset, fields); err != nil {
			return fmt.Errorf("prepare source %T error: %w", src, err)
		}
	}

	if err := fset.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("parse flags error: %w", err)
	}

	return nil
}

func (c *cliConfigs) Parse(ctx context.Context, srcs []sources.Source) error {
	for i := len(srcs) - 1; i >= 0; i-- {
		src := srcs[i]
		if err := src.Parse(ctx); err != nil {
			return fmt.Errorf("parse config from source %T error: %w", src, err)
		}
	}

	for name, handler := range c.configs {
		if err := handler.Callback(ctx); err != nil {
			return fmt.Errorf("init config %q error: %w", name, err)
		}
	}

	return nil
}
