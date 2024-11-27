/*
package clic is short for CLI Config. It implements a framework to load/parse cli configuration from a file, environment or flags.

Usage:

  - Register functions ([RegisterWithCallback] and [RegisterAndGet]) should be called in `init()` of a package, or before [Init] calls.
  - [Init] function should be called at the beginning of `main()`, before calling functions in other sub-packages.
  - [Init] function must not be called in `init()`, because other sub-packages may not be initialized at that time.
*/
package clic

import "context"

/*
RegisterWithCallback registers a  `Config` struct with the `name` as the scope name and `callback` function to consume the parsed instance of `Config`. [Init] function parses configuration from a file, environment or flags, stores into an instance of `Config`, then call `callback` with that instance. Then `callback` function could initialize global instances. If `callback` returns an error, [Init] fails and returns a wrapped error.

Example:

	package log

	import (
		"fmt"
		"log/slog"
		"os"

		"github.com/googollee/clic"
	)

	var logger *slog.Logger

	type config struct {
		Format string `clic:"format,json,the format of logging [json,text]"`
	}

	func initLogger(ctx context.Context, cfg *config) error {
		switch cfg.Format{
		case "json":
			logger = slog.New(slog.NewJSONHandler(os.Stderr))
		case "text":
			logger = slog.New(slog.NewTextHandler(os.Stderr))
		default:
			return fmt.Errorf("invalid log format: %q", cfg.Format)
	}
		return nil
	}

	func init() {
		clic.RegisterWithCallback("log", initLogger)
	}
*/
func RegisterWithCallback[Config any](name string, callback func(ctx context.Context, cfg *Config) error) {
}

/*
RegisterAndGet registers a `Config` struct with the `name` as the scope name and returns a function `getter` to get the `Config` instance from the context. [Init] function parses configuration from a file, environment or flags, stores a `Config` instance into the returned context.

Example:

	package main

	import (
		"library/database"
		"github.com/googollee/clic"
	)

	func main() {
		dbConfig := clic.RegisterAndGet[database.Config]("database")

		ctx, err := clic.Init(context.Background())
		if err != nil {
			panic(err)
		}

		db := database.New(dbConfig(ctx))
	}
*/
func RegisterAndGet[Config any](name string) (getter func(ctx context.Context) *Config) {
	return
}

// Init parses configuration from a file, environment or flags. It returns a new context which could be used to retreive values registered with [RegisterAndGet].
func Init(ctx context.Context) (context.Context, error) {
	return nil, nil
}
