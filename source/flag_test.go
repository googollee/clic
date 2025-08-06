package source

import (
	"bytes"
	"context"
	"flag"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFlag(t *testing.T) {
	tests := []struct {
		name                   string
		options                []FlagOption
		wantHelp               string
		args                   []string
		wantA1, wantA2, wantA3 string
	}{
		{
			name:     "FromValue",
			options:  []FlagOption{},
			wantHelp: "  -a1 value\n    \ta1 (default a1)\n  -l1.a2 value\n    \ta2 (default a2)\n  -l2.l3.a3 value\n    \ta3 (default a3)\n",
			args:     []string{"-a1", "123", "-l1.a2", "abc", "-l2.l3.a3", "xyz"},
			wantA1:   "123",
			wantA2:   "abc",
			wantA3:   "xyz",
		},
		{
			name:     "FromDefault",
			options:  []FlagOption{},
			wantHelp: "  -a1 value\n    \ta1 (default a1)\n  -l1.a2 value\n    \ta2 (default a2)\n  -l2.l3.a3 value\n    \ta3 (default a3)\n",
			args:     []string{},
			wantA1:   "a1",
			wantA2:   "a2",
			wantA3:   "a3",
		},
		{
			name:     "WithPrefixSplitter",
			options:  []FlagOption{FlagPrefix("demo"), FlagSplitter("_")},
			wantHelp: "  -demo_a1 value\n    \ta1 (default a1)\n  -demo_l1_a2 value\n    \ta2 (default a2)\n  -demo_l2_l3_a3 value\n    \ta3 (default a3)\n",
			args:     []string{"-demo_a1", "123", "-demo_l1_a2", "abc", "-demo_l2_l3_a3", "xyz"},
			wantA1:   "123",
			wantA2:   "abc",
			wantA3:   "xyz",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := flag.NewFlagSet("", flag.ContinueOnError)

			src := Flag(fset, tc.options...)
			a1, a2, a3 = "a1", "a2", "a3"

			if err := src.Register(fields); err != nil {
				t.Fatalf("src.Prepare(fields) returns error: %v", err)
			}

			var output bytes.Buffer
			fset.SetOutput(&output)
			fset.PrintDefaults()

			if diff := cmp.Diff(output.String(), tc.wantHelp); diff != "" {
				t.Errorf("output diff: (-got, +want)\n%s", diff)
			}

			if err := src.Parse(context.Background(), tc.args); err != nil {
				t.Fatalf("src.Parse() should return no error, which is not: %v", err)
			}

			if got, want := a1, tc.wantA1; got != want {
				t.Errorf("after src.Parse(), a1 = %q, want: %q", got, want)
			}
			if got, want := a2, tc.wantA2; got != want {
				t.Errorf("after src.Parse(), a2 = %q, want: %q", got, want)
			}
			if got, want := a3, tc.wantA3; got != want {
				t.Errorf("after src.Parse(), a3 = %q, want: %q", got, want)
			}
		})
	}

	t.Run("InvalidValue", func(t *testing.T) {
		fset := flag.NewFlagSet("", flag.ContinueOnError)
		src := Flag(fset)

		var output bytes.Buffer
		fset.SetOutput(&output)
		t.Logf("output:\n%s", output.String())

		if err := src.Register(fields); err != nil {
			t.Fatalf("src.Prepare(fields) returns error: %v", err)
		}

		args := []string{"-flag.a", "123"}

		if err := src.Parse(context.Background(), args); err == nil {
			t.Fatalf("src.Parse() = %v, want an error", err)
		}
	})
}

func TestFlagError(t *testing.T) {
	tests := []struct {
		name    string
		options []FlagOption
	}{
		{"EmptySplitter", []FlagOption{FlagSplitter("")}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := flag.NewFlagSet("", flag.ContinueOnError)
			src := Flag(fset, tc.options...)
			err := src.Error()
			if err == nil {
				t.Errorf("src().Error() want an error, which is not")
			}

			if errRegister := src.Register(fields); err != errRegister {
				t.Errorf("src().Prepare() = %v, src().Error() = %v, they should be same", errRegister, err)
			}

			if errParse := src.Parse(context.Background(), []string{}); err != errParse {
				t.Errorf("src().Parse() = %v, src().Error() = %v, they should be same", errParse, err)
			}
		})
	}
}
