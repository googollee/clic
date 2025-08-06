/*
package clic is short for CLI Config. It implements a framework to load/parse cli configuration from a file, environment or flags.

Usage:
  - Register functions ([RegisterWithCallback] and [RegisterAndGet]) should be called in "init()" of a package, or before [Init] calls.
  - [Init] function should be called at the beginning of "main()", before calling functions in other sub-packages.
  - [Init] function must not be called in "init()", because other sub-packages may not be initialized at that time.
*/
package clic

import (
	"context"
	"flag"
	"fmt"
	"os"
)

var CommandLine = NewSet(flag.CommandLine)

/*
Register registers a "Config" value with the "name" as the scope name. The value is filled after calling [Parse] function.

Example:

	package main

	func main() {
		ctx := context.Background()

		var dbCfg database.Config
		clic.Register("database", &db)

		clic.Parse(ctx)

		db := database.New(dbCfg)
	}
*/
func Register(prefix string, value any) error {
	return CommandLine.RegisterValue(prefix, value)
}

/*
RegisterCallback registers a callback function with the "name" as the scope name. The callback is called after calling [Parse] function.

Example:

	package main

	type Log struct {
		Level slog.Level `clic:"level,INFO,the level of log"`
	}

	func initLogLevel(ctx context.Context, cfg *Log) {
		slog.SetLevel(cfg.Level)
	}

	func main() {
		ctx := context.Background()

		clic.RegisterCallback("log", initLogLevel)

		clic.Parse(ctx)
	}
*/
func RegisterCallback(prefix string, f any) error {
	return CommandLine.RegisterCallback(prefix, f)
}

// Parse parses configuration from [DefaultSources] and [os.Args].
//
// If any error happens during calling, "Parse()" prints that error on Stderr and calls [os.Exit] to exit with "125" code.
func Parse(ctx context.Context) {
	if err := CommandLine.Parse(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "parse config error:", err)
		os.Exit(125)
	}
}
