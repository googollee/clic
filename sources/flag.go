package sources

import "flag"

type FlagOption func(*FlagSource) error

type FlagSource struct {
	set           *flag.FlagSet
	prefix        string
	splitter      string
	withHelp      bool
	withGenerator fileGenerator
}

func Flag(opt ...FlagOption) Source {
	return nil
}