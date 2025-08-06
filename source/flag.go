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

type flagSource struct {
	splitter string

	fset FlagSet
	err  error
}

func Flag(opt ...FlagOption) Source {
	ret := flagSource{
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

func (s *flagSource) Register(fset FlagSet, fields []structtags.Field) error {
	if s.err != nil {
		return s.err
	}

	s.fset = fset

	for _, field := range fields {
		names := field.Name

		key := strings.ToLower(strings.Join(names, s.splitter))
		fset.TextVar(&field, key, field, field.Description)
	}

	return nil
}

func (s *flagSource) Parse(ctx context.Context, args []string) error {
	if s.err != nil {
		return s.err
	}

	if err := s.fset.Parse(args); err != nil {
		return err
	}

	return nil
}
