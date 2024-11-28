package structtags

import (
	"fmt"
	"reflect"
	"strings"
)

type Field struct {
	Name          []string
	DefaultString string
	Description   string
	Parser        ParseFieldFunc
	Value         reflect.Value
}

func (f Field) String() string {
	return fmt.Sprintf("%v", f.Value.Interface())
}

func (f Field) Set(str string) error {
	return f.Parser(f.Value, str)
}

func (f Field) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f Field) UnmarshalText(buf []byte) error {
	return f.Parser(f.Value, string(buf))
}

func ParseStruct(v reflect.Value, name []string) ([]Field, error) {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("can't parse a %s type (names: %v), should be a struct", v.Kind(), name)
	}

	var ret []Field

	t := v.Type()
	vfields := reflect.VisibleFields(t)

	for _, vfield := range vfields {
		if vfield.Anonymous {
			continue
		}

		f := getFieldTag(vfield, name)

		vfieldValue := v.FieldByIndex(vfield.Index)
		vfieldType := vfieldValue.Type()

		if vfieldValue.Kind() == reflect.Pointer {
			vfieldType = vfieldType.Elem()

			if vfieldValue.IsNil() {
				vfieldValue.Set(reflect.New(vfieldType))
			}

			vfieldValue = vfieldValue.Elem()
		}

		if parser := getParseFieldFunc(vfieldType); parser != nil {
			f.Parser = parser
			f.Value = vfieldValue
			if f.DefaultString != "" {
				if err := f.Parser(f.Value, f.DefaultString); err != nil {
					return nil, fmt.Errorf("can't parse default value %q for field %v: %w", f.DefaultString, f.Name, err)
				}
			}
			ret = append(ret, f)
			continue
		}

		fields, err := ParseStruct(vfieldValue, f.Name)
		if err != nil {
			return nil, err
		}
		ret = append(ret, fields...)
	}

	return ret, nil
}

func getFieldTag(sfield reflect.StructField, name []string) (ret Field) {
	tagStr := sfield.Tag.Get("clic")
	tagArray := strings.SplitN(tagStr, ",", 3)
	name = name[:]

	switch len(tagArray) {
	case 3:
		ret.Name = append(name, tagArray[0])
		ret.DefaultString = tagArray[1]
		ret.Description = tagArray[2]
		return
	case 2:
		ret.Name = append(name, tagArray[0])
		ret.DefaultString = tagArray[1]
		return
	case 1:
		if tagStr != "" {
			ret.Name = append(name, tagArray[0])
			return
		}
	}

	if sfield.Name != "" {
		ret.Name = append(name, sfield.Name)
	}

	return
}
