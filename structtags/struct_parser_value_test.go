package structtags

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type testValueStruct struct {
	Int   int  `clic:"int,10"`
	PInt  *int `clic:"pint,20"`
	Inner struct {
		Dur time.Duration `clic:"dur,1h"`
	} `clic:"inner"`
}

func TestParseStructValue(t *testing.T) {
	var value testValueStruct
	fields, err := ParseStruct(reflect.ValueOf(&value).Elem(), []string{"test"})
	if err != nil {
		t.Fatalf("ParseStruct(%T) returns an error: %v, want no error", value, err)
	}

	t.Run("DefaultValue", func(t *testing.T) {
		i := 20
		wantDefault := testValueStruct{
			Int:  10,
			PInt: &i,
		}
		wantDefault.Inner.Dur = time.Hour

		if diff := cmp.Diff(value, wantDefault); diff != "" {
			t.Errorf("Diff: (-got, +want)\n%s", diff)
		}
	})

	fieldStrings := []string{"20", "40", "2h10m20s"}
	i := 40
	want := testValueStruct{
		Int:  20,
		PInt: &i,
	}
	want.Inner.Dur = 2*time.Hour + 10*time.Minute + 20*time.Second

	t.Run("TextUnmarshal", func(t *testing.T) {
		if got, want := len(fields), len(fieldStrings); got != want {
			t.Fatalf("len(fields) = %d should equal to len(fieldStrings) = %d, which is not", got, want)
		}

		for i, fieldString := range fieldStrings {
			if err := fields[i].TextUnmarshal([]byte(fieldString)); err != nil {
				t.Fatalf("Field %v: Parse(%q) returns an error: %v, want no error", fields[i].Name, fieldString, err)
			}

			if got, want := fields[i].String(), fieldString; got != want {
				t.Errorf("Field %v: String() = %q, want: %q", fields[i].Name, got, want)
			}
		}

		if diff := cmp.Diff(value, want); diff != "" {
			t.Errorf("Diff: (-got, +want)\n%s", diff)
		}
	})

	t.Run("ValueSet", func(t *testing.T) {
		if got, want := len(fields), len(fieldStrings); got != want {
			t.Fatalf("len(fields) = %d should equal to len(fieldStrings) = %d, which is not", got, want)
		}

		for i, fieldString := range fieldStrings {
			if err := fields[i].Set(fieldString); err != nil {
				t.Fatalf("Field %v: Parse(%q) returns an error: %v, want no error", fields[i].Name, fieldString, err)
			}

			if got, want := fields[i].String(), fieldString; got != want {
				t.Errorf("Field %v: String() = %q, want: %q", fields[i].Name, got, want)
			}
		}

		if diff := cmp.Diff(value, want); diff != "" {
			t.Errorf("Diff: (-got, +want)\n%s", diff)
		}
	})
}

type testInvalidValueStruct struct {
	Int int `clic:"int,abc"`
}

func TestParseStructInvalidValue(t *testing.T) {
	var value testInvalidValueStruct
	_, err := ParseStruct(reflect.ValueOf(&value).Elem(), []string{"test"})
	if err == nil {
		t.Fatalf("ParseStruct(%T) returns no error, want an parsing error", value)
	}
}
