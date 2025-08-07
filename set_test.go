package clic_test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"testing"
	"unsafe"

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

func ExampleSet_showHelp() {
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
	fset := flag.NewFlagSet("", flag.ContinueOnError)
	var helpOutput bytes.Buffer
	fset.SetOutput(&helpOutput)
	set := clic.NewSet(fset,
		// The order means srouce priority, flag > config file > env
		source.Flag(source.FlagSplitter(".")),
		source.File(source.FilePathFlag("config"), source.FileFormat(source.JSON{})),
		source.Env(source.EnvSplitter("_")),
	)

	// args
	args := []string{"-h"}

	// main code
	ctx := context.Background()

	var db Database
	set.RegisterValue("database", &db)
	set.RegisterCallback("log", initLog)

	if err := set.Parse(ctx, args); !errors.Is(err, flag.ErrHelp) {
		log.Fatal("parse error:", err)
	}

	fmt.Println(strings.ReplaceAll(helpOutput.String(), "\t", "    "))
	// Output:
	// Usage:
	//   -config string
	//         the path of the config file
	//   -database.driver value
	//         the driver of the database [sqlite3,mysql,postgres] (default sqlite3)
	//   -database.url value
	//         the url of the database (default ./database.sqlite)
	//   -log.level value
	//         the level of the log [DEBUG,INFO,WARN,ERROR] (default INFO)
}

func TestInvalidRegister(t *testing.T) {
	type C struct{ Int int }
	var c C

	type Invalid struct {
		Int int `clic:"int,invalid,int"`
	}
	var invalid Invalid

	var vInterface any
	var vFunc func()
	var vChan chan struct{}
	var vUnsafePointer unsafe.Pointer
	var vMap map[string]string

	tests := []struct {
		name  string
		value any
	}{
		{"value", &c},
		{"callback", &c},

		{"nil_value", nil},
		{"non_ptr", "abc"},
		{"invalid_default", &invalid},
		{"p_interface", &vInterface},
		{"p_func", &vFunc},
		{"p_chan", &vChan},
		{"p_unsafe_pointer", &vUnsafePointer},
		{"p_map", &vMap},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := flag.NewFlagSet("", flag.ContinueOnError)
			set := clic.NewSet(fset)

			set.RegisterValue("value", &c)
			set.RegisterCallback("callback", func(context.Context, *C) error { return nil })

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
	type C struct{ Int int }
	var c C

	type Invalid struct {
		Int int `clic:"int,invalid,int"`
	}

	tests := []struct {
		name string
		fn   any
	}{
		{"value", func(context.Context, *C) error { return nil }},
		{"callback", func(context.Context, *C) error { return nil }},

		{"nil_func", nil},
		{"non_func", "str"},
		{"invalid_default", func(context.Context, *Invalid) error { return nil }},
		{"non_ptr", func(context.Context, int) error { return nil }},
		{"non_ctx", func(int, *int) error { return nil }},
		{"non_return", func(context.Context, *int) {}},
		{"non_error", func(context.Context, *int) int { return 1 }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := flag.NewFlagSet("", flag.ContinueOnError)
			set := clic.NewSet(fset)

			set.RegisterValue("value", &c)
			set.RegisterCallback("callback", func(context.Context, *C) error { return nil })

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
	wantErr := errors.New("error")

	type C struct {
		Int int `clic:"int,0,int"`
	}
	initC := func(context.Context, *C) error {
		return wantErr
	}

	fset := flag.NewFlagSet("", flag.ContinueOnError)
	set := clic.NewSet(fset)

	set.RegisterCallback("callback", initC)

	ctx := context.Background()
	if err := set.Parse(ctx, []string{}); !errors.Is(err, wantErr) {
		t.Errorf("set.Parse() = %v, want: %v", err, wantErr)
	}
}

func TestInvalidValue(t *testing.T) {
	type C struct {
		Int int `clic:"int,0,int"`
	}
	var c C

	fset := flag.NewFlagSet("", flag.ContinueOnError)
	set := clic.NewSet(fset)

	set.RegisterValue("value", &c)

	t.Setenv("VALUE_INT", "invalid_int")

	ctx := context.Background()
	if err := set.Parse(ctx, []string{}); err == nil {
		t.Errorf("set.Parse() = %v, want an error", err)
	}
}
