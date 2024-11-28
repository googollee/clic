package structtags

import (
	"encoding"
	"flag"
	"fmt"
	"reflect"
	"slices"
	"testing"
	"time"
)

func testFieldParser[T comparable](t *testing.T, str string, want T, wantOk bool) {
	var got T
	tv := reflect.TypeOf(got)
	parser := getParseFieldFunc(tv)
	if parser == nil {
		t.Fatalf("can't find a parser for type %q.", tv)
	}

	err := parser(reflect.ValueOf(&got).Elem(), str)
	gotOk := err == nil
	if gotOk != wantOk {
		whatWant := "want an error"
		if wantOk {
			whatWant = "want no error"
		}
		t.Errorf("parser[%s](%q) = error(%v), %s", tv, str, err, whatWant)

		return
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("parser[%s](%q) returns (%T)%v, want: (%T)%v", tv, str, got, got, want, want)
	}
}

func TestParseFieldString(t *testing.T) {
	succeeded := []struct {
		input string
	}{
		{"a"},
		{"alongstring"},
	}

	for _, tc := range succeeded {
		testFieldParser(t, tc.input, tc.input, true)
	}
}

func TestParseFieldInt(t *testing.T) {
	succeeded := []struct {
		input string
		want  int64
	}{
		{"0", 0},
		{"1", 1},
		{"9", 9},
		{"09", 9},
		{"10", 10},
		{"99", 99},
		{"-0", 0},
		{"-1", -1},
		{"-9", -9},
		{"-10", -10},
		{"-99", -99},
	}

	failure := []struct {
		input string
	}{
		{"abc"},
		{""},
		{"123a"},
		{"-123a"},
	}

	for _, tc := range succeeded {
		testFieldParser(t, tc.input, int(tc.want), true)
		testFieldParser(t, tc.input, int64(tc.want), true)
		testFieldParser(t, tc.input, int32(tc.want), true)
		testFieldParser(t, tc.input, int16(tc.want), true)
		testFieldParser(t, tc.input, int8(tc.want), true)
	}

	for _, tc := range failure {
		testFieldParser(t, tc.input, int(0), false)
		testFieldParser(t, tc.input, int64(0), false)
		testFieldParser(t, tc.input, int32(0), false)
		testFieldParser(t, tc.input, int16(0), false)
		testFieldParser(t, tc.input, int8(0), false)
	}
}

func TestParseFieldUint(t *testing.T) {
	succeeded := []struct {
		input string
		want  uint64
	}{
		{"0", 0},
		{"1", 1},
		{"9", 9},
		{"10", 10},
		{"99", 99},
		{"0b1100", 0b1100},
		{"0xFF", 0xFF},
		{"0o123", 0o123},
		{"0123", 0o123},
	}

	failure := []struct {
		input string
	}{
		{"abc"},
		{""},
		{"123a"},
		{"-123a"},
		{"-1"},
		{"0b2"},
		{"0xG"},
		{"0o9"},
		{"09"},
	}

	for _, tc := range succeeded {
		testFieldParser(t, tc.input, uint(tc.want), true)
		testFieldParser(t, tc.input, uint64(tc.want), true)
		testFieldParser(t, tc.input, uint32(tc.want), true)
		testFieldParser(t, tc.input, uint16(tc.want), true)
		testFieldParser(t, tc.input, uint8(tc.want), true)
	}

	for _, tc := range failure {
		testFieldParser(t, tc.input, uint(0), false)
		testFieldParser(t, tc.input, uint64(0), false)
		testFieldParser(t, tc.input, uint32(0), false)
		testFieldParser(t, tc.input, uint16(0), false)
		testFieldParser(t, tc.input, uint8(0), false)
	}
}

func TestParseFieldFloat(t *testing.T) {
	succeeded := []struct {
		input string
		want  float64
	}{
		{"0.1", 0.1},
		{"1.1", 1.1},
		{"1", 1},
		{"-0.1", -0.1},
		{"-1.1", -1.1},
		{"-1", -1},
	}

	failure := []struct {
		input string
	}{
		{"abc"},
		{""},
		{"123a"},
		{"-123a"},
		{"0b2"},
		{"0xG"},
		{"0o9"},
	}

	for _, tc := range succeeded {
		testFieldParser(t, tc.input, float64(tc.want), true)
		testFieldParser(t, tc.input, float32(tc.want), true)
	}

	for _, tc := range failure {
		testFieldParser(t, tc.input, float64(0), false)
		testFieldParser(t, tc.input, float32(0), false)
	}
}

func TestParseFieldBool(t *testing.T) {
	succeeded := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
		{"True", true},
		{"False", false},
		{"1", true},
		{"0", false},
	}

	failure := []struct {
		input string
	}{
		{"abc"},
		{""},
		{"123a"},
		{"-123a"},
		{"-1"},
		{"0b2"},
		{"0xG"},
		{"0o9"},
		{"09"},
	}

	for _, tc := range succeeded {
		testFieldParser(t, tc.input, tc.want, true)
	}

	for _, tc := range failure {
		testFieldParser(t, tc.input, bool(false), false)
	}
}

func TestParseFieldTextUnmarshaler(t *testing.T) {
	var _ encoding.TextUnmarshaler = &time.Time{}

	succeeded := []struct {
		input string
		want  time.Time
	}{
		{"2024-11-25T20:50:00Z", time.Date(2024, 11, 25, 20, 50, 0, 0, time.UTC)},
	}

	failure := []struct {
		input string
	}{
		{"123"},
	}

	for _, tc := range succeeded {
		testFieldParser(t, tc.input, tc.want, true)
	}

	for _, tc := range failure {
		testFieldParser(t, tc.input, time.Time{}, false)
	}
}

type enumValue string

func (v enumValue) String() string {
	return string(v)
}

func (v *enumValue) Set(s string) error {
	if !slices.Contains([]string{"a", "b", "c"}, s) {
		return fmt.Errorf("invalid enumValue: %q", s)
	}

	*v = enumValue(s)
	return nil
}

func TestParseFieldFlagValue(t *testing.T) {
	var value enumValue
	var _ flag.Value = &value

	succeeded := []struct {
		input string
		want  enumValue
	}{
		{"a", enumValue("a")},
	}

	failure := []struct {
		input string
	}{
		{"123"},
	}

	for _, tc := range succeeded {
		testFieldParser(t, tc.input, tc.want, true)
	}

	for _, tc := range failure {
		testFieldParser(t, tc.input, enumValue(""), false)
	}
}

func TestParseFieldDuration(t *testing.T) {
	succeeded := []struct {
		input string
		want  time.Duration
	}{
		{"2h45m", time.Hour*2 + time.Minute*45},
	}

	failure := []struct {
		input string
	}{
		{"123"},
		{"abc"},
	}

	for _, tc := range succeeded {
		testFieldParser(t, tc.input, tc.want, true)
	}

	for _, tc := range failure {
		testFieldParser(t, tc.input, time.Duration(0), false)
	}
}
