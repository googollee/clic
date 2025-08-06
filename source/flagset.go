package source

import "encoding"

type FlagSet interface {
	PrintDefaults()
	Parse([]string) error
	Parsed() bool

	TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string)
	StringVar(v *string, name string, defaultValue string, usage string)
	BoolVar(v *bool, name string, defaultValue bool, usage string)
}
