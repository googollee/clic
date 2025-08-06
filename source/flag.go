package source

import (
	"context"
	"errors"
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
	fset     FlagSet

	err error
}

func Flag(fset FlagSet, opt ...FlagOption) Source {
	ret := flagSource{
		prefix:   "",
		splitter: ".",
		fset:     fset,
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

func (s *flagSource) Register(fields []structtags.Field) error {
	if s.err != nil {
		return s.err
	}

	for _, field := range fields {
		names := field.Name
		if s.prefix != "" {
			names = append([]string{s.prefix}, field.Name...)
		}

		key := strings.ToLower(strings.Join(names, s.splitter))
		s.fset.TextVar(&field, key, field, field.Description)
	}

	return nil
}

func (s *flagSource) Parse(ctx context.Context, args []string) error {
	if s.err != nil {
		return s.err
	}

	if !s.fset.Parsed() {
		if err := s.fset.Parse(args); err != nil {
			return err
		}
	}

	return nil
}
