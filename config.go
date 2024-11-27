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
RegisterWithCallback registers an Initializer instance `value` with the `name` as the scope name. [Init] function parses configuration from a file, environment or flags, stores into `value`, then call `value.Init()`. The `value` provider could initialize global instances in `value.Init()`.

Example:

	package subpack

	var singletonInstance int

	type config struct {
		Int int `clic:"int,10,the value of the global instance"`
	}

	func (c config) Init(ctx context.Context) error {
		singletonInstance = c.Int
	}

	func init() {
		clic.RegisterWithCallback[config]("subpack")
	}
*/
func RegisterWithCallback[T any](name string, callback func(ctx context.Context, cfg *T) error) {}

/*
RegisterAndGet registers an instance `value` with the `name` as the scope name and returns a function to get parsed `value` from the context. [Init] function parses configuration from a file, environment or flags, stores into `value`, then call `value.Init()`.

Example:

	package main

	func main() {
		dbConfig := clic.RegisterAndGet[db.Config]("database")

		ctx, err := clic.Init(context.Background())
		if err != nil {
			panic(err)
		}

		db := database.New(dbConfig(ctx))
	}
*/
func RegisterAndGet[T any](name string) (getter func(ctx context.Context) *T) {
	return
}

// Init parses configuration from a file, environment or flags. It returns a new context which could be used to retreive values registered with [RegisterAndGet].
func Init(ctx context.Context) (context.Context, error) {
	return nil, nil
}
