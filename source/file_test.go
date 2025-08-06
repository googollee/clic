package source

import (
	"bytes"
	"flag"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFile(t *testing.T) {
	tests := []struct {
		name                   string
		options                []FileOption
		wantHelp               string
		args                   []string
		wantA1, wantA2, wantA3 string
	}{
		{
			name:     "FromValue",
			options:  []FileOption{},
			wantHelp: "  -config string\n    \tthe path of the config file\n",
			args:     []string{"-config", "./testdata/valid.json"},
			wantA1:   "123",
			wantA2:   "abc",
			wantA3:   "xyz",
		},
		{
			name:     "FromDefault",
			options:  []FileOption{},
			wantHelp: "  -config string\n    \tthe path of the config file\n",
			args:     []string{"-config", "./testdata/empty.json"},
			wantA1:   "a1",
			wantA2:   "a2",
			wantA3:   "a3",
		},
		{
			name:     "WithFlagFormat",
			options:  []FileOption{FileFormat(JSON{}), FilePathFlag("c")},
			wantHelp: "  -c string\n    \tthe path of the config file\n",
			args:     []string{"-c", "./testdata/valid.json"},
			wantA1:   "123",
			wantA2:   "abc",
			wantA3:   "xyz",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := flag.NewFlagSet("", flag.ContinueOnError)
			src := File(tc.options...)
			a1, a2, a3 = "a1", "a2", "a3"

			if err := src.Register(fset, fields); err != nil {
				t.Fatalf("src.Prepare(fields) returns error: %v", err)
			}

			var output bytes.Buffer
			fset.SetOutput(&output)
			fset.PrintDefaults()

			if diff := cmp.Diff(output.String(), tc.wantHelp); diff != "" {
				t.Errorf("output diff: (-got, +want)\n%s", diff)
			}

			if err := fset.Parse(tc.args); err != nil {
				t.Fatalf("fset.Parse() error: %v", err)
			}

			if err := src.Parse(t.Context(), tc.args); err != nil {
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

	t.Run("InvalidFile", func(t *testing.T) {
		fset := flag.NewFlagSet("", flag.ContinueOnError)
		src := File()

		var output bytes.Buffer
		fset.SetOutput(&output)
		t.Logf("output:\n%s", output.String())

		if err := src.Register(fset, fields); err != nil {
			t.Fatalf("src.Prepare(fields) returns error: %v", err)
		}

		args := []string{"-config", "testdata/not_exist.json"}
		if err := fset.Parse(args); err != nil {
			t.Fatalf("fset.Parse() error: %v", err)
		}

		if err := src.Parse(t.Context(), args); err == nil {
			t.Fatalf("src.Parse() error: %v, want an error", err)
		}
	})

	t.Run("EmptyConfigFile", func(t *testing.T) {
		fset := flag.NewFlagSet("", flag.ContinueOnError)
		src := File()

		if err := src.Register(fset, fields); err != nil {
			t.Fatalf("src.Prepare(fields) returns error: %v", err)
		}

		args := []string{"-config", ""}
		if err := fset.Parse(args); err != nil {
			t.Fatalf("fset.Parse() error: %v", err)
		}

		if err := src.Parse(t.Context(), args); err != nil {
			t.Fatalf("src.Parse() = %v, want no error", err)
		}
	})

	t.Run("InvalidFlag", func(t *testing.T) {
		fset := flag.NewFlagSet("", flag.ContinueOnError)
		src := File()

		if err := src.Register(fset, fields); err != nil {
			t.Fatalf("src.Prepare(fields) returns error: %v", err)
		}

		args := []string{"-nonexist_flag", ""}
		if err := fset.Parse(args); err == nil {
			t.Fatalf("fset.Parse() = nil, want an error")
		}

		if err := src.Parse(t.Context(), args); err != nil {
			t.Fatalf("src.Parse() = %v, want no error", err)
		}
	})
}

func TestFileError(t *testing.T) {
	tests := []struct {
		name    string
		options []FileOption
	}{
		{"EmptyCodec", []FileOption{FileFormat(nil)}},
		{"EmptyPathFlag", []FileOption{FilePathFlag("")}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := flag.NewFlagSet("", flag.ContinueOnError)
			src := File(tc.options...)
			err := src.Error()
			if err == nil {
				t.Errorf("src().Error() want an error, which is not")
			}

			errPrepare := src.Register(fset, fields)
			if err != errPrepare {
				t.Errorf("src().Prepare() = %v, src().Error() = %v, they should be same", errPrepare, err)
			}

			args := []string{}
			if err := fset.Parse(args); err != nil {
				t.Fatalf("fset.Parse() error: %v", err)
			}

			errParse := src.Parse(t.Context(), args)
			if err != errPrepare {
				t.Errorf("src().Parse() = %v, src().Error() = %v, they should be same", errParse, err)
			}
		})
	}
}
