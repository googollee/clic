package sources

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/googollee/clic/structtags"
)

func TestEnv(t *testing.T) {
	tests := []struct {
		name                   string
		options                []EnvOption
		envs                   map[string]string
		wantA1, wantA2, wantA3 string
	}{
		{
			name:    "FromValue",
			options: []EnvOption{},
			envs: map[string]string{
				"CLIC_A1":       "123",
				"CLIC_L1_A2":    "abc",
				"CLIC_L2_L3_A3": "xyz",
			},
			wantA1: "123",
			wantA2: "abc",
			wantA3: "xyz",
		},
		{
			name:    "FromDefault",
			options: []EnvOption{},
			envs:    map[string]string{},
			wantA1:  "a1",
			wantA2:  "a2",
			wantA3:  "a3",
		},
		{
			name:    "WithPrefixSplitter",
			options: []EnvOption{EnvPrefix("DEMO"), EnvSplitter("__")},
			envs: map[string]string{
				"DEMO__A1":         "123",
				"DEMO__L1__A2":     "abc",
				"DEMO__L2__L3__A3": "xyz",
			},
			wantA1: "123",
			wantA2: "abc",
			wantA3: "xyz",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			src := Env(tc.options...)
			a1, a2, a3 = "a1", "a2", "a3"

			if err := src.Prepare(nil, fields); err != nil {
				t.Fatalf("src.Prepare(fields) returns error: %v", err)
			}

			for key, value := range tc.envs {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("os.Setenv(%q, %q) returns error: %v", key, value, err)
				}
			}
			defer func() {
				for key := range tc.envs {
					os.Unsetenv(key)
				}
			}()

			if err := src.Parse(context.Background()); err != nil {
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
		flagSet := flag.NewFlagSet("", flag.ContinueOnError)
		src := File()

		var output bytes.Buffer
		flagSet.SetOutput(&output)
		t.Logf("output:\n%s", output.String())

		if err := src.Prepare(flagSet, fields); err != nil {
			t.Fatalf("src.Prepare(fields) returns error: %v", err)
		}

		args := []string{"-config", "testdata/not_exist.json"}
		if err := flagSet.Parse(args); err != nil {
			t.Fatalf("flagSet.Parse() error: %v", err)
		}

		if err := src.Parse(context.Background()); err == nil {
			t.Fatalf("src.Parse() error: %v, want an error", err)
		}
	})
}

func TestEnvOptionError(t *testing.T) {
	tests := []struct {
		name    string
		options []EnvOption
	}{
		{"EmptySplitter", []EnvOption{EnvSplitter("")}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			src := Env(tc.options...)
			err := src.Error()
			if err == nil {
				t.Errorf("src().Error() want an error, which is not")
			}

			errPrepare := src.Prepare(nil, nil)
			if err != errPrepare {
				t.Errorf("src().Prepare() = %v, src().Error() = %v, they should be same", errPrepare, err)
			}

			errParse := src.Parse(context.Background())
			if err != errPrepare {
				t.Errorf("src().Parse() = %v, src().Error() = %v, they should be same", errParse, err)
			}
		})
	}
}

func parserWithError(v reflect.Value, str string) error {
	return fmt.Errorf("error!")
}

func TestEnvValueError(t *testing.T) {
	var i int
	tests := []struct {
		fields []structtags.Field
		envs   map[string]string
	}{
		{
			fields: []structtags.Field{
				{Name: []string{"int"}, Description: "int", Parser: parserWithError, Value: reflect.ValueOf(&i).Elem()},
			},
			envs: map[string]string{"CLIC_INT": "abc"}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.envs), func(t *testing.T) {
			i = 0
			src := Env()

			if err := src.Prepare(nil, tc.fields); err != nil {
				t.Errorf("src().Prepare() = %v, should be no error", err)
			}

			for key, value := range tc.envs {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("os.Setenv(%q, %q) returns error: %v", key, value, err)
				}
			}
			defer func() {
				for key := range tc.envs {
					os.Unsetenv(key)
				}
			}()

			if err := src.Parse(context.Background()); err == nil {
				t.Errorf("src().Parse() = %v, should be an error", err)
			}
		})
	}
}
