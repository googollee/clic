package sources

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/googollee/clic/structtags"
)

var ErrQuitEarly = errors.New("quit early")

type FlagOption func(*flagSource) error

func FlagSplitter(splitter string) FlagOption {
	return func(s *flagSource) error {
		if splitter == "" {
			return fmt.Errorf("invalid flag splitter: %q", splitter)
		}
		s.splitter = splitter
		return nil
	}
}

func FlagPrefix(prefix string) FlagOption {
	return func(s *flagSource) error {
		s.prefix = prefix
		return nil
	}
}

type flagSource struct {
	prefix   string
	splitter string

	err error
}

func Flag(opt ...FlagOption) Source {
	ret := flagSource{
		prefix:   "",
		splitter: ".",
	}

	for _, opt := range opt {
		if err := opt(&ret); err != nil {
			ret.err = err
		}
	}

	return &ret
}

func (s *flagSource) Error() error {
	return s.err
}

func (s *flagSource) Prepare(fset *flag.FlagSet, fields []structtags.Field) error {
	if s.err != nil {
		return s.err
	}

	for _, field := range fields {
		field := field
		names := field.Name
		if s.prefix != "" {
			names = append([]string{s.prefix}, field.Name...)
		}
		fset.TextVar(&field, strings.Join(names, s.splitter), field, field.Description)
	}

	return nil
}

func (s *flagSource) Parse(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}

	return nil
}
