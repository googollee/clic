package structtags

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type parseFieldFunc func(v reflect.Value, str string) error

// getParseFieldFunc doesn't accept a Pointer Type. Caller should check pointer and pass in Pointer.Elem().
func getParseFieldFunc(t reflect.Type) parseFieldFunc {
	if reflect.PointerTo(t).Implements(umarshalerType) {
		return parseFieldUmarshaler
	}

	if t == durationType {
		return parseFieldDuration
	}

	return parsers[t]
}

var parsers = map[reflect.Type]parseFieldFunc{
	reflect.TypeOf(""):          parseFieldString,
	reflect.TypeOf(int(0)):      parseFieldInt[int],
	reflect.TypeOf(int64(0)):    parseFieldInt[int64],
	reflect.TypeOf(int32(0)):    parseFieldInt[int32],
	reflect.TypeOf(int16(0)):    parseFieldInt[int16],
	reflect.TypeOf(int8(0)):     parseFieldInt[int8],
	reflect.TypeOf(uint(0)):     parseFieldUint[uint],
	reflect.TypeOf(uint64(0)):   parseFieldUint[uint64],
	reflect.TypeOf(uint32(0)):   parseFieldUint[uint32],
	reflect.TypeOf(uint16(0)):   parseFieldUint[uint16],
	reflect.TypeOf(uint8(0)):    parseFieldUint[uint8],
	reflect.TypeOf(float64(0)):  parseFieldFloat[float64],
	reflect.TypeOf(float32(0)):  parseFieldFloat[float32],
	reflect.TypeOf(bool(false)): parseFieldBool,
}

func parseFieldString(v reflect.Value, str string) error {
	v.Set(reflect.ValueOf(str))
	return nil
}

func parseFieldInt[Int ~int | ~int8 | ~int16 | ~int32 | ~int64](v reflect.Value, str string) error {
	i64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse %q to an integer: %w", str, err)
	}

	v.Set(reflect.ValueOf(Int(i64)))
	return nil
}

func parseFieldUint[UInt ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](v reflect.Value, str string) error {
	base := 10
	if strings.HasPrefix(str, "0x") {
		base = 16
		str = str[2:]
	} else if strings.HasPrefix(str, "0b") {
		base = 2
		str = str[2:]
	} else if strings.HasPrefix(str, "0o") {
		base = 8
		str = str[2:]
	} else if strings.HasPrefix(str, "0") && str != "0" {
		base = 8
		str = str[1:]
	}

	u64, err := strconv.ParseUint(str, base, 64)
	if err != nil {
		return fmt.Errorf("can't parse %q to an unsigned-integer: %w", str, err)
	}

	v.Set(reflect.ValueOf(UInt(u64)))
	return nil
}

func parseFieldFloat[Float ~float32 | ~float64](v reflect.Value, str string) error {
	f64, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("can't parse %q to a float: %w", str, err)
	}

	v.Set(reflect.ValueOf(Float(f64)))
	return nil
}

func parseFieldBool(v reflect.Value, str string) error {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return fmt.Errorf("can't parse %q to a bool: %w", str, err)
	}

	v.Set(reflect.ValueOf(b))
	return nil
}

var umarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()

func parseFieldUmarshaler(v reflect.Value, str string) error {
	if v.Kind() == reflect.Struct {
		v = v.Addr()
	}

	return v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(str))
}

var durationType = reflect.TypeFor[time.Duration]()

func parseFieldDuration(v reflect.Value, str string) error {
	dur, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(dur))
	return nil
}
