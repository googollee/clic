package structtags

import (
	"fmt"
	"reflect"
	"strings"
)

type Field struct {
	Name          string
	DefaultString string
	Description   string
	Parser        parseFunc
	Index         []int
	Children      []Field
}

func Parse(t reflect.Type) ([]Field, error) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	vfields := reflect.VisibleFields(t)
	ret := make([]Field, 0, len(vfields))
	for _, vfield := range vfields {
		var f Field

		tagStr := vfield.Tag.Get("clic")
		tagArray := strings.SplitN(tagStr, ",", 3)

		switch len(tagArray) {
		case 3:
			f.Name = tagArray[0]
			f.DefaultString = tagArray[1]
			f.Description = tagArray[2]
		case 2:
			f.Name = tagArray[0]
			f.DefaultString = tagArray[1]
		case 1:
			f.Name = tagArray[0]
		default:
			f.Name = vfield.Name
		}

		f.Index = vfield.Index

		vfieldType := vfield.Type
		if vfieldType.Kind() == reflect.Pointer {
			vfieldType = vfieldType.Elem()
		}

		switch vfieldType.Kind() {
		case reflect.Struct:
			var err error
			f.Children, err = Parse(vfieldType)
			if err != nil {
				return nil, err
			}
		default:
			parser, ok := parsers[vfieldType]
			if !ok {
				return nil, fmt.Errorf("can't parse type %s.", vfieldType)
			}
			f.Parser = parser
		}
	}

	return ret, nil
}
