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
	"os"

	"github.com/googollee/clic/sources"
	"github.com/googollee/clic/structtags"
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
	if _, exist := handlers[name]; exist {
		panic("already registered a config with name " + name)
	}

	handler := config[Config]{}
	handlers[name] = &handler
	getter = handler.Get

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
	if _, exist := handlers[name]; exist {
		panic("already registered a config with name " + name)
	}

	handler := config[Config]{
		callback: callback,
	}
	handlers[name] = &handler
}

// Init parses configuration from a file, environment or flags.
//
// If any error happens during calling, "Init()" prints that error on Stderr and calls [os.Exit] to exit with "125" code.
func Init(ctx context.Context) {
	srcs := sources.Default

	fset := flag.CommandLine
	var showHelp bool
	fset.BoolVar(&showHelp, "help", false, "show the usage")
	fset.BoolVar(&showHelp, "h", false, "show the usage")

	var fields []structtags.Field
	for name, handler := range handlers {
		f, err := structtags.ParseStruct(handler.Value(), []string{name})
		if err != nil {
			panic("parse config of " + name + " error: " + err.Error())
		}
		fields = append(fields, f...)
	}

	for i := 0; i < len(srcs); i++ {
		src := srcs[i]
		if err := src.Prepare(fset, fields); err != nil {
			panic("prepare error: " + err.Error())
		}
	}

	if err := fset.Parse(os.Args[1:]); err != nil {
		panic("parse flags error: " + err.Error())
	}

	if showHelp {
		fset.PrintDefaults()
		os.Exit(125)
		return
	}

	for i := len(srcs) - 1; i >= 0; i-- {
		src := srcs[i]
		if err := src.Parse(ctx); err != nil {
			panic("parse config error: " + err.Error())
		}
	}

	for name, handler := range handlers {
		if err := handler.Callback(ctx); err != nil {
			panic("init config " + name + " error " + err.Error())
		}
	}
}
