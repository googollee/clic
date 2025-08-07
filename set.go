package clic

import (
	"context"
	"fmt"

	"github.com/googollee/clic/source"
	"github.com/googollee/clic/structtags"
)

var DefaultSources = []source.Source{
	source.Flag(source.FlagSplitter(".")),
	source.File(source.FilePathFlag("config"), source.FileFormat(source.JSON{})),
	source.Env(source.EnvSplitter("_")),
}

type Set struct {
	fset    source.FlagSet
	sources []source.Source
	configs map[string]*config
	fields  []structtags.Field
}

func NewSet(fset source.FlagSet, source ...source.Source) *Set {
	if len(source) == 0 {
		source = DefaultSources
	}

	return &Set{
		fset:    fset,
		sources: source,
		configs: make(map[string]*config),
	}
}

func (s *Set) RegisterValue(prefix string, value any) {
	if err := s.register(prefix, newConfigValue(value)); err != nil {
		panic(err)
	}
}

func (s *Set) RegisterCallback(prefix string, callback any) {
	if err := s.register(prefix, newConfigCallback(callback)); err != nil {
		panic(err)
	}
}

func (s *Set) Parse(ctx context.Context, args []string) error {
	for i := range len(s.sources) {
		src := s.sources[i]
		if err := src.Register(s.fset, s.fields); err != nil {
			return fmt.Errorf("prepare source %T error: %w", src, err)
		}
	}

	if s.fset != nil && !s.fset.Parsed() {
		if err := s.fset.Parse(args); err != nil {
			return err
		}
	}

	for i := len(s.sources) - 1; i >= 0; i-- {
		src := s.sources[i]
		if err := src.Parse(ctx, args); err != nil {
			return fmt.Errorf("parse config from source %T error: %w", src, err)
		}
	}

	for name, handler := range s.configs {
		if err := handler.Callback(ctx); err != nil {
			return fmt.Errorf("init config %q error: %w", name, err)
		}
	}

	return nil
}

func (s *Set) register(prefix string, config *config) error {
	fields, err := structtags.ParseStruct(config.Value(), []string{prefix})
	if err != nil {
		return err
	}

	if _, exist := s.configs[prefix]; exist {
		return fmt.Errorf("already registered a config with prefix %s", prefix)
	}

	s.fields = append(s.fields, fields...)
	s.configs[prefix] = config

	return nil
}
