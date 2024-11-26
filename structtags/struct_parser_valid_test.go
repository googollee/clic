package structtags

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestStruct struct {
	Int  int     `clic:"int,10,integer value"`
	PInt *int    `clic:"pint,20,pointer to integer"`
	Str  string  `clic:"str,some str,string value"`
	PStr *string `clic:"pstr,some str,pointer to string"`
}

type TestEmbedStruct struct {
	TestStruct
}

var testStructFields = []Field{
	{
		Name:          []string{"test", "int"},
		DefaultString: "10",
		Description:   "integer value",
	},
	{
		Name:          []string{"test", "pint"},
		DefaultString: "20",
		Description:   "pointer to integer",
	},
	{
		Name:          []string{"test", "str"},
		DefaultString: "some str",
		Description:   "string value",
	},
	{
		Name:          []string{"test", "pstr"},
		DefaultString: "some str",
		Description:   "pointer to string",
	},
}

type TestPartClic struct {
	NoDesc   string `clic:"no_desc,default"`
	OnlyName string `clic:"only_name"`
	NoTag    string
}

var testPartClicFields = []Field{
	{
		Name:          []string{"test", "no_desc"},
		DefaultString: "default",
		Description:   "",
	},
	{
		Name:          []string{"test", "only_name"},
		DefaultString: "",
		Description:   "",
	},
	{
		Name:          []string{"test", "NoTag"},
		DefaultString: "",
		Description:   "",
	},
}

type TestLayersStruct struct {
	Struct       TestStruct       `clic:"struct"`
	PStruct      *TestStruct      `clic:"pstruct"`
	EmbedStruct  TestEmbedStruct  `clic:"embed"`
	PEmbedStruct *TestEmbedStruct `clic:"pembed"`
}

var testLayerStructFields = []Field{
	{
		Name:          []string{"test", "struct", "int"},
		DefaultString: "10",
		Description:   "integer value",
	},
	{
		Name:          []string{"test", "struct", "pint"},
		DefaultString: "20",
		Description:   "pointer to integer",
	},
	{
		Name:          []string{"test", "struct", "str"},
		DefaultString: "some str",
		Description:   "string value",
	},
	{
		Name:          []string{"test", "struct", "pstr"},
		DefaultString: "some str",
		Description:   "pointer to string",
	},
	{
		Name:          []string{"test", "pstruct", "int"},
		DefaultString: "10",
		Description:   "integer value",
	},
	{
		Name:          []string{"test", "pstruct", "pint"},
		DefaultString: "20",
		Description:   "pointer to integer",
	},
	{
		Name:          []string{"test", "pstruct", "str"},
		DefaultString: "some str",
		Description:   "string value",
	},
	{
		Name:          []string{"test", "pstruct", "pstr"},
		DefaultString: "some str",
		Description:   "pointer to string",
	},
	{
		Name:          []string{"test", "embed", "int"},
		DefaultString: "10",
		Description:   "integer value",
	},
	{
		Name:          []string{"test", "embed", "pint"},
		DefaultString: "20",
		Description:   "pointer to integer",
	},
	{
		Name:          []string{"test", "embed", "str"},
		DefaultString: "some str",
		Description:   "string value",
	},
	{
		Name:          []string{"test", "embed", "pstr"},
		DefaultString: "some str",
		Description:   "pointer to string",
	},
	{
		Name:          []string{"test", "pembed", "int"},
		DefaultString: "10",
		Description:   "integer value",
	},
	{
		Name:          []string{"test", "pembed", "pint"},
		DefaultString: "20",
		Description:   "pointer to integer",
	},
	{
		Name:          []string{"test", "pembed", "str"},
		DefaultString: "some str",
		Description:   "string value",
	},
	{
		Name:          []string{"test", "pembed", "pstr"},
		DefaultString: "some str",
		Description:   "pointer to string",
	},
}

func compareField(x, y Field) bool {
	if !reflect.DeepEqual(x.Name, y.Name) {
		return false
	}
	if x.DefaultString != y.DefaultString {
		return false
	}
	if x.Description != y.Description {
		return false
	}
	return true
}

func TestParseStruct(t *testing.T) {
	var partValue TestPartClic
	var structValue TestStruct
	var embedValue TestEmbedStruct
	var layerValue TestLayersStruct

	tests := []struct {
		value      any
		wantFields []Field
	}{
		{&partValue, testPartClicFields},
		{&structValue, testStructFields},
		{&embedValue, testStructFields},
		{&layerValue, testLayerStructFields},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%T", tc.value), func(t *testing.T) {
			for _, value := range []reflect.Value{reflect.ValueOf(tc.value).Elem(), reflect.ValueOf(tc.value)} {
				got, err := ParseStruct(value, []string{"test"})
				if err != nil {
					t.Errorf("ParseStruct(%v, ['test']) returns an error: %v, want: no error", value.Type(), err)
				}

				if diff := cmp.Diff(got, tc.wantFields, cmp.Comparer(compareField)); diff != "" {
					t.Errorf("ParseStruct(%v, ['test']) diff: (-got, +want)\n:%s", value.Type(), diff)
				}
			}
		})
	}
}
