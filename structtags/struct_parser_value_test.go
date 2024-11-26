package structtags

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type testValueStruct struct {
	Int   int  `clic:"int"`
	PInt  *int `clic:"pint"`
	Inner struct {
		Dur time.Duration `clic:"dur"`
	} `clic:"inner"`
}

func TestParseStructValue(t *testing.T) {
	var value testValueStruct
	fields, err := ParseStruct(reflect.ValueOf(&value).Elem(), []string{"test"})
	if err != nil {
		t.Fatalf("ParseStruct(%T) returns an error: %v, want no error", value, err)
	}

	fieldStrings := []string{"10", "20", "1h"}
	i := 20
	want := testValueStruct{
		Int:  10,
		PInt: &i,
	}
	want.Inner.Dur = time.Hour

	if got, want := len(fields), len(fieldStrings); got != want {
		t.Fatalf("len(fields) = %d should equal to len(fieldStrings) = %d, which is not", got, want)
	}

	for i, fieldString := range fieldStrings {
		if err := fields[i].Parser(fields[i].Value, fieldString); err != nil {
			t.Fatalf("Field %v: Parse(%q) returns an error: %v, want no error", fields[i].Name, fieldString, err)
		}
	}

	if diff := cmp.Diff(value, want); diff != "" {
		t.Errorf("Diff: (-got, +want)\n%s", diff)
	}
}
