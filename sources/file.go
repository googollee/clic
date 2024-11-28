package sources

import (
	"context"
	"flag"
	"fmt"
	"reflect"

	"github.com/googollee/clic/structtags"
)

type FileCodec interface {
	TagName() string
	ExtName() string
	Encode(path string, v any) error
	Decode(path string, v any) error
}

type FileOption func(*fileSource) error

func FileFormat(codec FileCodec) FileOption {
	return func(s *fileSource) error {
		if codec == nil {
			return fmt.Errorf("invalid codec: %v", codec)
		}
		s.codec = codec
		return nil
	}
}

func FileFlagSet(set *flag.FlagSet) FileOption {
	return func(s *fileSource) error {
		if set == nil {
			return fmt.Errorf("invalid flag set: %v", set)
		}
		s.flagSet = set
		return nil
	}
}

type fileSource struct {
	codec        FileCodec
	flagSet      *flag.FlagSet
	filepathFlag string
	filepath     string
	err          error

	value reflect.Value
}

func File(options ...FileOption) Source {
	ret := fileSource{
		codec:        JSON{},
		flagSet:      flag.CommandLine,
		filepathFlag: "config",
	}
	ret.filepath = "./config." + ret.codec.ExtName()

	for _, option := range options {
		if err := option(&ret); err != nil {
			ret.err = err
		}
	}

	return &ret
}

func (s *fileSource) Error() error {
	return s.err
}

func (s *fileSource) Prepare(fields []structtags.Field) error {
	if s.err != nil {
		return s.err
	}

	s.value = newFromFields(fields, 0, s.codec.TagName()+":\"%s\"")
	s.flagSet.StringVar(&s.filepath, s.filepathFlag, s.filepath, "the path of the config file")

	return nil
}

func (s *fileSource) Parse(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}

	return s.codec.Decode(s.filepath, s.value.Interface())
}
