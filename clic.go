/*
package clic is short for CLI Config. It implements a framework to load/parse cli configuration from a file, environment or flags.

Usage:
  - Register functions ([Register] and [RegisterCallback]) should be called in `func init()` in packages, or before calling [Parse].
  - [Parse] function should be called at the beginning of "main()", before calling functions in other packages.
  - [Parse] function must not be called in `func init()`, because other sub-packages may not finish initialization at that time.

See examples for the usage.
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
func Register(prefix string, value any) {
	CommandLine.RegisterValue(prefix, value)
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
func RegisterCallback(prefix string, f any) {
	CommandLine.RegisterCallback(prefix, f)
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
