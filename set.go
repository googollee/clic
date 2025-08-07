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
	fset     source.FlagSet
	sources  []source.Source
	configs  map[string]*config
	withHelp bool
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
	adapter := newConfigValue(value)

	if _, exist := s.configs[prefix]; exist {
		panic(fmt.Errorf("already registered a config with prefix %s", prefix))
	}

	s.configs[prefix] = adapter
}

func (s *Set) RegisterCallback(prefix string, callback any) {
	adapter := newConfigCallback(callback)

	if _, exist := s.configs[prefix]; exist {
		panic(fmt.Errorf("already registered a config with prefix %s", prefix))
	}

	s.configs[prefix] = adapter
}

func (s *Set) Parse(ctx context.Context, args []string) error {
	if err := s.prepare(); err != nil {
		return err
	}

	if s.fset != nil {
		if !s.fset.Parsed() {
			if err := s.fset.Parse(args); err != nil {
				return err
			}
		}
	}

	if err := s.parse(ctx, args); err != nil {
		return err
	}

	return nil
}

func (s *Set) prepare() error {
	var fields []structtags.Field
	for name, handler := range s.configs {
		f, err := structtags.ParseStruct(handler.Value(), []string{name})
		if err != nil {
			return fmt.Errorf("parse config %q error: %w", name, err)
		}
		fields = append(fields, f...)
	}

	for i := range len(s.sources) {
		src := s.sources[i]
		if err := src.Register(s.fset, fields); err != nil {
			return fmt.Errorf("prepare source %T error: %w", src, err)
		}
	}

	return nil
}

func (s *Set) parse(ctx context.Context, args []string) error {
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
