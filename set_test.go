package clic_test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		{"non_func", "str"},
		{"non_ptr", func(context.Context, int) error { return nil }},
		{"non_ctx", func(int, *int) error { return nil }},
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

func TestCallbackError(t *testing.T) {
	fset := flag.NewFlagSet("", flag.PanicOnError)
	set := clic.NewSet(fset)

	wantErr := fmt.Errorf("error")
	set.RegisterCallback("callback", func(context.Context, *struct{ Int int }) error {
		return wantErr
	})

	ctx := context.Background()
	if err := set.Parse(ctx, []string{"-callback.int", "1"}); !errors.Is(err, wantErr) {
		t.Errorf("set.Parse() == %v, want %v", err, wantErr)
	}
}

func TestHelp(t *testing.T) {
	var output bytes.Buffer

	fset := flag.NewFlagSet("", flag.ContinueOnError)
	fset.SetOutput(&output)

	set := clic.NewSet(fset)

	set.RegisterCallback("callback", func(context.Context, *struct{ Int int }) error {
		return nil
	})

	ctx := context.Background()
	if got, want := set.Parse(ctx, []string{"-h"}), flag.ErrHelp; !errors.Is(got, want) {
		t.Errorf("set.Parse() == %v, want %s", got, want)
	}

	t.Logf("output: %q", output.String())
	wantOutput := "Usage:\n  -callback.int value\n    \t (default 0)\n  -config string\n    \tthe path of the config file\n"
	if diff := cmp.Diff(output.String(), wantOutput); diff != "" {
		t.Errorf("output diff(-got, +want):\n%s", diff)
	}

}
