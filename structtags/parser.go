package structtags

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type parseFunc func(v reflect.Value, str string) error

// getParseFunc doesn't accept a Pointer Type. Caller should check pointer and pass in Pointer.Elem().
func getParseFunc(t reflect.Type) parseFunc {
	if reflect.PointerTo(t).Implements(umarshalerType) {
		return parseUmarshaler
	}

	if t == durationType {
		return parseDuration
	}

	return parsers[t]
}

var parsers = map[reflect.Type]parseFunc{
	reflect.TypeOf(""):          parseString,
	reflect.TypeOf(int(0)):      parseInt[int],
	reflect.TypeOf(int64(0)):    parseInt[int64],
	reflect.TypeOf(int32(0)):    parseInt[int32],
	reflect.TypeOf(int16(0)):    parseInt[int16],
	reflect.TypeOf(int8(0)):     parseInt[int8],
	reflect.TypeOf(uint(0)):     parseUint[uint],
	reflect.TypeOf(uint64(0)):   parseUint[uint64],
	reflect.TypeOf(uint32(0)):   parseUint[uint32],
	reflect.TypeOf(uint16(0)):   parseUint[uint16],
	reflect.TypeOf(uint8(0)):    parseUint[uint8],
	reflect.TypeOf(float64(0)):  parseFloat[float64],
	reflect.TypeOf(float32(0)):  parseFloat[float32],
	reflect.TypeOf(bool(false)): parseBool,
}

func parseString(v reflect.Value, str string) error {
	v.Set(reflect.ValueOf(str))
	return nil
}

func parseInt[Int ~int | ~int8 | ~int16 | ~int32 | ~int64](v reflect.Value, str string) error {
	i64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse %q to an integer: %w", str, err)
	}

	v.Set(reflect.ValueOf(Int(i64)))
	return nil
}

func parseUint[UInt ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](v reflect.Value, str string) error {
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

func parseFloat[Float ~float32 | ~float64](v reflect.Value, str string) error {
	f64, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("can't parse %q to a float: %w", str, err)
	}

	v.Set(reflect.ValueOf(Float(f64)))
	return nil
}

func parseBool(v reflect.Value, str string) error {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return fmt.Errorf("can't parse %q to a bool: %w", str, err)
	}

	v.Set(reflect.ValueOf(b))
	return nil
}

var umarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()

func parseUmarshaler(v reflect.Value, str string) error {
	if v.Kind() == reflect.Struct {
		v = v.Addr()
	}

	return v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(str))
}

var durationType = reflect.TypeFor[time.Duration]()

func parseDuration(v reflect.Value, str string) error {
	dur, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(dur))
	return nil
}
