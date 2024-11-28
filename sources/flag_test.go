package sources

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googollee/clic/structtags"
)

func TestFlagSource(t *testing.T) {
	var a1, a2, a3 string
	fields := []structtags.Field{
		{Name: []string{"a1"}, Description: "a1", Parser: parserString, Value: reflect.ValueOf(&a1).Elem()},
		{Name: []string{"l1", "a2"}, Description: "a2", Parser: parserString, Value: reflect.ValueOf(&a2).Elem()},
		{Name: []string{"l2", "l3", "a3"}, Description: "a3", Parser: parserString, Value: reflect.ValueOf(&a3).Elem()},
	}

	t.Run("Help", func(t *testing.T) {
		flagSet := flag.NewFlagSet("", flag.ContinueOnError)
		src := Flag(FlagSet(flagSet), FlagWithHelp(true))

		if err := src.Prepare(fields); err != nil {
			t.Fatalf("src.Prepare(fields) returns error: %v", err)
		}

		var output bytes.Buffer
		flagSet.SetOutput(&output)

		os.Args = []string{"-help"}

		ctx := context.Background()
		if err := src.Parse(ctx); !errors.Is(err, ErrQuitEarly) {
			t.Fatalf("src.Parse() should return ErrQuitEarly, which is not: %v", err)
		}

		want := "  -flag.a1 value\n    \ta1\n  -flag.l1.a2 value\n    \ta2\n  -flag.l2.l3.a3 value\n    \ta3\n  -help\n    \tShow the usage\n"
		if diff := cmp.Diff(output.String(), want); diff != "" {
			t.Errorf("output diff: (-got, +want)\n%s", diff)
		}
	})

	t.Run("HelpPrefixSplitter", func(t *testing.T) {
		flagSet := flag.NewFlagSet("", flag.ContinueOnError)
		src := Flag(FlagSet(flagSet), FlagWithHelp(true), FlagPrefix("demo"), FlagSplitter("_"))

		if err := src.Prepare(fields); err != nil {
			t.Fatalf("src.Prepare(fields) returns error: %v", err)
		}

		var output bytes.Buffer
		flagSet.SetOutput(&output)

		os.Args = []string{"-help"}

		ctx := context.Background()
		if err := src.Parse(ctx); !errors.Is(err, ErrQuitEarly) {
			t.Fatalf("src.Parse() should return ErrQuitEarly, which is not: %v", err)
		}

		want := "  -demo_a1 value\n    \ta1\n  -demo_l1_a2 value\n    \ta2\n  -demo_l2_l3_a3 value\n    \ta3\n  -help\n    \tShow the usage\n"
		if diff := cmp.Diff(output.String(), want); diff != "" {
			t.Errorf("output diff: (-got, +want)\n%s", diff)
		}
	})

	t.Run("Value", func(t *testing.T) {
		flagSet := flag.NewFlagSet("", flag.ContinueOnError)
		src := Flag(FlagSet(flagSet))

		if err := src.Prepare(fields); err != nil {
			t.Fatalf("src.Prepare(fields) returns error: %v", err)
		}

		os.Args = []string{"-flag.a1", "123", "-flag.l1.a2", "abc", "-flag.l2.l3.a3", "xyz"}

		ctx := context.Background()
		if err := src.Parse(ctx); err != nil {
			t.Fatalf("src.Parse() should return no error, which is not: %v", err)
		}

		if got, want := a1, "123"; got != want {
			t.Errorf("after src.Parse(), a1 = %q, want: %q", got, want)
		}
		if got, want := a2, "abc"; got != want {
			t.Errorf("after src.Parse(), a2 = %q, want: %q", got, want)
		}
		if got, want := a3, "xyz"; got != want {
			t.Errorf("after src.Parse(), a3 = %q, want: %q", got, want)
		}
	})

	t.Run("ValuePrefixSplitter", func(t *testing.T) {
		flagSet := flag.NewFlagSet("", flag.ContinueOnError)
		src := Flag(FlagSet(flagSet), FlagPrefix("demo"), FlagSplitter("_"))

		if err := src.Prepare(fields); err != nil {
			t.Fatalf("src.Prepare(fields) returns error: %v", err)
		}

		os.Args = []string{"-demo_a1", "123", "-demo_l1_a2", "abc", "-demo_l2_l3_a3", "xyz"}

		ctx := context.Background()
		if err := src.Parse(ctx); err != nil {
			t.Fatalf("src.Parse() should return no error, which is not: %v", err)
		}

		if got, want := a1, "123"; got != want {
			t.Errorf("after src.Parse(), a1 = %q, want: %q", got, want)
		}
		if got, want := a2, "abc"; got != want {
			t.Errorf("after src.Parse(), a2 = %q, want: %q", got, want)
		}
		if got, want := a3, "xyz"; got != want {
			t.Errorf("after src.Parse(), a3 = %q, want: %q", got, want)
		}
	})

	t.Run("InvalidValue", func(t *testing.T) {
		flagSet := flag.NewFlagSet("", flag.ContinueOnError)
		src := Flag(FlagSet(flagSet))

		if err := src.Prepare(fields); err != nil {
			t.Fatalf("src.Prepare(fields) returns error: %v", err)
		}

		os.Args = []string{"-flag.a", "123"}

		ctx := context.Background()
		if err := src.Parse(ctx); err == nil {
			t.Fatalf("src.Parse() should return an error, which is not")
		}
	})
}

func TestFlagSourceError(t *testing.T) {
	tests := []struct {
		name    string
		options []FlagOption
	}{
		{"EmptyPrefix", []FlagOption{FlagPrefix("")}},
		{"EmpytSplitter", []FlagOption{FlagSplitter("")}},
		{"NilSet", []FlagOption{FlagSet(nil)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			src := Flag(tc.options...)
			err := src.Error()
			if err == nil {
				t.Errorf("src().Error() want an error, which is not")
			}

			errPrepare := src.Prepare(nil)
			if err != errPrepare {
				t.Errorf("src().Prepare() = %v, src().Error() = %v, they should be same", errPrepare, err)
			}

			errParse := src.Parse(nil)
			if err != errPrepare {
				t.Errorf("src().Parse() = %v, src().Error() = %v, they should be same", errParse, err)
			}
		})
	}
}
