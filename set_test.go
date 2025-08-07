package clic_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"testing"

	"github.com/googollee/clic"
	"github.com/googollee/clic/source"
)

func ExampleSet() {
	// structs
	type Database struct {
		Driver string `clic:"driver,sqlite3,the driver of the database [sqlite3,mysql,postgres]"`
		URL    string `clic:"url,./database.sqlite,the url of the database"`
	}
	type Log struct {
		Level slog.Level `clic:"level,INFO,the level of the log [DEBUG,INFO,WARN,ERROR]"`
	}
	initLog := func(ctx context.Context, log *Log) error {
		fmt.Println("set log level:", log.Level)
		return nil
	}

	// set with sources
	fset := flag.NewFlagSet("", flag.ExitOnError)
	set := clic.NewSet(fset,
		// The order means srouce priority, flag > config file > env
		source.Flag(source.FlagSplitter(".")),
		source.File(source.FilePathFlag("config"), source.FileFormat(source.JSON{})),
		source.Env(source.EnvSplitter("_")),
	)

	// args
	args := []string{"-log.level", "WARN", "-database.driver", "driver", "-database.url", "url", "other_cmd"}

	// main code
	ctx := context.Background()

	var db Database
	set.RegisterValue("database", &db)
	set.RegisterCallback("log", initLog)

	if err := set.Parse(ctx, args); err != nil {
		log.Fatal("parse error:", err)
	}

	fmt.Println("database:", db)
	fmt.Println("remain args:", fset.Args())

	// Output:
	// set log level: WARN
	// database: {driver url}
	// remain args: [other_cmd]
}

func TestInvalidRegister(t *testing.T) {
	var value int
	tests := []struct {
		name  string
		value any
	}{
		{"value", &value},
		{"callback", &value},

		{"nil_value", nil},
		{"non_ptr", "abc"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := flag.NewFlagSet("", flag.ContinueOnError)
			set := clic.NewSet(fset)

			var i int
			set.RegisterValue("value", &i)
			set.RegisterCallback("callback", func(context.Context, *int) error { return nil })

			defer func() {
				if r := recover(); r == nil {
					t.Error("set.Register() passes, want a panic")
				}
			}()

			set.RegisterValue(tc.name, tc.value)
		})
	}
}

func TestInvalidRegisterCallback(t *testing.T) {
	tests := []struct {
		name string
		fn   any
	}{
		{"value", func(context.Context, *int) error { return nil }},
		{"callback", func(context.Context, *int) error { return nil }},

		{"nil_func", nil},
		{"non_ptr", func(context.Context, int) error { return nil }},
		{"non_ctx", func(*int) error { return nil }},
		{"non_return", func(context.Context, *int) {}},
		{"non_error", func(context.Context, *int) int { return 1 }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := flag.NewFlagSet("", flag.ContinueOnError)
			set := clic.NewSet(fset)

			var i int
			set.RegisterValue("value", &i)
			set.RegisterCallback("callback", func(context.Context, *int) error { return nil })

			defer func() {
				if r := recover(); r == nil {
					t.Error("set.RegisterCallback() passes, want a panic")
				}
			}()

			set.RegisterCallback(tc.name, tc.fn)
		})
	}
}
