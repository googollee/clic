/*
package clic is short for CLI Config. It implements a framework to load/parse cli configuration from a file, environment or flags.

Usage:
  - Register functions ([RegisterWithCallback] and [RegisterAndGet]) should be called in "init()" of a package, or before [Init] calls.
  - [Init] function should be called at the beginning of "main()", before calling functions in other sub-packages.
  - [Init] function must not be called in "init()", because other sub-packages may not be initialized at that time.
*/
package clic

import "context"

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
		// `clic` tag format: <name>,<default value>,<description>
		Level string `clic:"level,info,the minimum level of logging: <debug|info|warn|error>"`
	}

	func initLogger(ctx context.Context, cfg *logConfig) error {
		switch cfg.Level {
		case "debug":
			slog.SetLogLoggerLevel(slog.LevelDebug)
		case "info":
			slog.SetLogLoggerLevel(slog.LevelInfo)
		case "warn":
			slog.SetLogLoggerLevel(slog.LevelWarn)
		case "error":
			slog.SetLogLoggerLevel(slog.LevelError)
		default:
			return fmt.Errorf("invalid log level: %q", cfg.Level)
		}

		return nil
	}
*/
func RegisterWithCallback[Config any](name string, callback func(ctx context.Context, cfg *Config) error) {
}

// Init parses configuration from a file, environment or flags.
//
// If any error happens during calling, "Init()" prints that error on Stderr and calls [os.Exit] to exit with "125" code.
func Init(ctx context.Context) {
}
