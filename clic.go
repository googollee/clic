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

	"github.com/googollee/clic/sources"
)

/*
RegisterAndGet registers a "Config" struct with the "name" as the scope name and returns a function "getter" to get the "Config" instance after calling [Init] function.

Example:

	package main

	import (
		// ...
	)

	var dbConfig = clic.RegisterAndGet[database.Config]("database")

	func main() {
		ctx := context.Background()
		clic.Init(ctx)

		db := database.New(dbConfig())
	}
*/
func RegisterAndGet[Config any](name string) (getter func() *Config) {
	adapter := config[Config]{}
	getter = adapter.Get

	if err := configs.Register(name, &adapter); err != nil {
		configs.ExitWithError(err)
	}

	return
}

/*
RegisterWithCallback registers a "Config" struct with the "name" as the scope name and a "callback" function to consume the parsed instance of "Config". [Init] function parses configuration from a file, environment or flags, stores into an instance of "Config", then calls "callback" with that instance. "callback" function could initialize global instances. If "callback" returns an error, [Init] fails and the binary exits.

Example:

	package main

	import (
		// ...
	)

	func main() {
		ctx := context.Background()
		clic.Init(ctx)

		slog.Info("Hello, clic!")
	}

	// To improve code organization and maintainability, consider moving the following code into a separate package.
	// If you do this, remember to import the new package into the "main" package to ensure proper functionality.

	func init() {
		clic.RegisterWithCallback(initLogger)
	}

	type logConfig struct {
		// [slog.Level] implements [encoding.TextMarshaler] and [encoding.TextUnmarshaler] to parse values from a string.
		// `clic` tag format: <name>,<default value>,<description>
		Level slog.Level `clic:"level,info,the minimum level of logging: <debug|info|warn|error>"`
	}

	func initLogger(ctx context.Context, cfg *logConfig) error {
		slog.SetLogLoggerLevel(cfg.Level)

		return nil
	}
*/
func RegisterWithCallback[Config any](name string, callback func(ctx context.Context, cfg *Config) error) {
	adapter := config[Config]{
		callback: callback,
	}
	if err := configs.Register(name, &adapter); err != nil {
		configs.ExitWithError(err)
	}
}

// Init parses configuration from a file, environment or flags.
//
// If any error happens during calling, "Init()" prints that error on Stderr and calls [os.Exit] to exit with "125" code.
func Init(ctx context.Context) {
	srcs := sources.Default

	fset := flag.CommandLine

	if err := configs.Prepare(srcs, fset); err != nil {
		configs.ExitWithError(err)
	}

	var help bool

	if configs.WithHelp() {
		fset.BoolVar(&help, "help", false, "show the usage")
		fset.BoolVar(&help, "h", false, "show the usage")
	}

	if err := fset.Parse(os.Args[1:]); err != nil {
		configs.ExitWithError(fmt.Errorf("parse flags error: %w", err))
	}

	if help {
		fset.PrintDefaults()
		os.Exit(0)
	}

	if err := configs.Parse(ctx, srcs); err != nil {
		configs.ExitWithError(err)
	}
}
