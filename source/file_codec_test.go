package source

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCodec(t *testing.T) {
	tests := []struct {
		codec   FileCodec
		content string
		wantTag string
		wantExt string
	}{
		{JSON{}, "{\"int\":1,\"str\":\"str\"}\n", "json", "json"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%T", tc.codec), func(t *testing.T) {
			if got, want := tc.codec.TagName(), tc.wantTag; got != want {
				t.Errorf("codec.TagName() = %q, want: %q", got, want)
			}

			if got, want := tc.codec.ExtName(), tc.wantExt; got != want {
				t.Errorf("codec.ExtName() = %q, want: %q", got, want)
			}

			fname := filepath.Join(t.TempDir(), "config."+tc.codec.ExtName())
			f, err := os.Create(fname)
			if err != nil {
				t.Fatalf("can't create temp file %q: %v", fname, err)
			}
			_, _ = f.WriteString(tc.content)
			f.Close()

			var value map[string]any
			if err := tc.codec.Decode(fname, &value); err != nil {
				t.Fatalf("tc.codec.Decode() returns error: %v", err)
			}

			if err := tc.codec.Encode(fname, &value); err != nil {
				t.Fatalf("tc.codec.Encode() returns error: %v", err)
			}

			gotContent, err := os.ReadFile(fname)
			if err != nil {
				t.Fatalf("os.ReadFile(%q) error: %v", fname, err)
			}

			if diff := cmp.Diff(string(gotContent), tc.content); diff != "" {
				t.Errorf("the diff content after decoding and encoding:\n%s", diff)
			}
		})

		t.Run(fmt.Sprintf("%TInvalidFile", tc.codec), func(t *testing.T) {
			i := 1
			fname := "/nonexist/file"

			if err := tc.codec.Decode(fname, &i); err == nil {
				t.Fatalf("tc.codec.Decode() should return an error, which is not")
			}

			if err := tc.codec.Encode(fname, &i); err == nil {
				t.Fatalf("tc.codec.Encode() should return an error, which is not")
			}
		})
	}
}
