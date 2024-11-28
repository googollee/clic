package sources

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/googollee/clic/structtags"
)

var ErrQuitEarly = errors.New("quit early")

type FlagOption func(*flagSource) error

func FlagSet(set *flag.FlagSet) FlagOption {
	return func(s *flagSource) error {
		if set == nil {
			return fmt.Errorf("invalid flag set: %v", set)
		}
		s.set = set
		return nil
	}
}

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
		if prefix == "" {
			return fmt.Errorf("invalid flag prefix: %q", prefix)
		}
		s.prefix = prefix
		return nil
	}
}

func FlagWithHelp(withHelp bool) FlagOption {
	return func(s *flagSource) error {
		s.withHelp = withHelp
		return nil
	}
}

type flagSource struct {
	set      *flag.FlagSet
	prefix   string
	splitter string
	withHelp bool

	err       error
	printHelp bool
}

func Flag(opt ...FlagOption) Source {
	ret := flagSource{
		set:      flag.CommandLine,
		prefix:   "flag",
		splitter: ".",
		withHelp: true,
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

func (s *flagSource) Prepare(fields []structtags.Field) error {
	if s.err != nil {
		return s.err
	}

	for _, field := range fields {
		field := field
		names := append([]string{s.prefix}, field.Name...)
		s.set.TextVar(&field, strings.Join(names, s.splitter), field, field.Description)
	}

	if s.withHelp {
		s.set.BoolVar(&s.printHelp, "help", s.printHelp, "Show the usage")
	}

	return nil
}

func (s *flagSource) Parse(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}

	err := s.set.Parse(os.Args)
	if err != nil {
		return fmt.Errorf("parse from flag error: %w", err)
	}

	if s.printHelp {
		s.set.PrintDefaults()
		return ErrQuitEarly
	}

	return nil
}
