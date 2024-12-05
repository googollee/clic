package sources

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/googollee/clic/structtags"
)

type EnvOption func(*envSource) error

func EnvPrefix(prefix string) EnvOption {
	return func(s *envSource) error {
		s.prefix = prefix
		return nil
	}
}

func EnvSplitter(splitter string) EnvOption {
	return func(s *envSource) error {
		if splitter == "" {
			return fmt.Errorf("invalid splitter: %q", splitter)
		}
		s.splitter = splitter
		return nil
	}
}

type envSource struct {
	prefix   string
	splitter string
	err      error
	fields   []structtags.Field
}

func Env(options ...EnvOption) Source {
	ret := envSource{
		prefix:   "clic",
		splitter: "_",
	}

	for _, option := range options {
		if err := option(&ret); err != nil {
			ret.err = err
		}
	}

	return &ret
}

func (s *envSource) Error() error {
	return s.err
}

func (s *envSource) Prepare(fset FlagSet, fields []structtags.Field) error {
	if s.err != nil {
		return s.err
	}

	s.fields = fields

	return nil
}

func (s *envSource) Parse(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}

	for _, field := range s.fields {
		name := field.Name
		if s.prefix != "" {
			name = append([]string{s.prefix}, field.Name...)
		}

		envKey := strings.ToUpper(strings.Join(name, s.splitter))
		envValue, exist := os.LookupEnv(envKey)
		if !exist {
			continue
		}

		if err := field.UnmarshalText([]byte(envValue)); err != nil {
			return fmt.Errorf("parse env (%s: %q) error: %w", envKey, envValue, err)
		}
	}

	return nil
}
