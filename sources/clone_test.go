package sources

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googollee/clic/structtags"
)

func TestFindEndFieldHaveSameIndex(t *testing.T) {
	tests := []struct {
		inputNames [][]string
		inputIndex int
		want       int
	}{
		{[][]string{}, 0, 0},
		{[][]string{{}}, 0, 0},

		{[][]string{{"a1"}, {"a2"}, {"a3"}}, 0, 1},
		{[][]string{{"a1"}, {"a2"}, {"a3"}}, 1, 0},

		{[][]string{{"a1", "b1"}, {"a1", "b2"}, {"a2"}}, 0, 2},
		{[][]string{{"a1", "b1"}, {"a1", "b2"}, {"a2"}}, 1, 1},

		{[][]string{{"a1", "b1"}, {"a1", "b2"}, {"a2", "b3"}}, 1, 1},
		{[][]string{{"a1", "b1"}, {"a1", "b1", "c1"}, {"a2", "b2"}}, 1, 2},

		{[][]string{{"a1", "b1"}, {"a2"}, {"a3", "b3"}}, 1, 1},

		{[][]string{{"a1", "b1"}, {"a2", "b2"}, {"a3", "b3"}}, 1, 1},

		{[][]string{{"a1"}, {"l1", "a2"}, {"l1", "l2", "a3"}}, 0, 1},
		{[][]string{{"l1", "a2"}, {"l1", "l2", "a3"}}, 0, 2},
	}

	for _, tc := range tests {
		fields := make([]structtags.Field, 0, len(tc.inputNames))
		for _, name := range tc.inputNames {
			fields = append(fields, structtags.Field{
				Name: name,
			})
		}

		got := findEndFieldHaveSameIndex(fields, tc.inputIndex)

		if got, want := got, tc.want; got != want {
			t.Errorf("findEndFieldHaveSameIndex(%v, %d) = %d, want: %d", tc.inputNames, tc.inputIndex, got, want)
		}
	}
}

func parserString(v reflect.Value, str string) error {
	v.Set(reflect.ValueOf(str))
	return nil
}

func TestNewFromFields(t *testing.T) {
	var a1, a2, a3 string
	tests := []struct {
		wantDefault string
		checkString string
		inputFields []structtags.Field
	}{
		{`{"a1":"layer1"}`, `{"a1":"123"}`, []structtags.Field{
			{Name: []string{"a1"}, Parser: parserString, Value: reflect.ValueOf(&a1).Elem()},
		}},
		{`{"l1":{"a2":"layer2"}}`, `{"l1":{"a2":"123"}}`, []structtags.Field{
			{Name: []string{"l1", "a2"}, Parser: parserString, Value: reflect.ValueOf(&a2).Elem()},
		}},
		{`{"l1":{"l2":{"a3":"layer3"}}}`, `{"l1":{"l2":{"a3":"123"}}}`, []structtags.Field{
			{Name: []string{"l1", "l2", "a3"}, Parser: parserString, Value: reflect.ValueOf(&a3).Elem()},
		}},

		{`{"a1":"layer1","l1":{"a2":"layer2"}}`, `{"a1":"abc","l1":{"a2":"123"}}`, []structtags.Field{
			{Name: []string{"a1"}, Parser: parserString, Value: reflect.ValueOf(&a1).Elem()},
			{Name: []string{"l1", "a2"}, Parser: parserString, Value: reflect.ValueOf(&a2).Elem()},
		}},

		{`{"a1":"layer1","l1":{"a2":"layer2","l2":{"a3":"layer3"}}}`, `{"a1":"123","l1":{"a2":"abc","l2":{"a3":"xyz"}}}`, []structtags.Field{
			{Name: []string{"a1"}, Parser: parserString, Value: reflect.ValueOf(&a1).Elem()},
			{Name: []string{"l1", "a2"}, Parser: parserString, Value: reflect.ValueOf(&a2).Elem()},
			{Name: []string{"l1", "l2", "a3"}, Parser: parserString, Value: reflect.ValueOf(&a3).Elem()},
		}},
		{`{"l1":{"a2":"layer2","l2":{"a3":"layer3"}},"a1":"layer1"}`, `{"l1":{"a2":"abc","l2":{"a3":"xyz"}},"a1":"123"}`, []structtags.Field{
			{Name: []string{"l1", "a2"}, Parser: parserString, Value: reflect.ValueOf(&a2).Elem()},
			{Name: []string{"l1", "l2", "a3"}, Parser: parserString, Value: reflect.ValueOf(&a3).Elem()},
			{Name: []string{"a1"}, Parser: parserString, Value: reflect.ValueOf(&a1).Elem()},
		}},
		{`{"a1":"layer1","l1":{"a2":"layer2"},"l2":{"l3":{"a3":"layer3"}}}`, `{"a1":"123","l1":{"a2":"abc"},"l2":{"l3":{"a3":"xyz"}}}`, []structtags.Field{
			{Name: []string{"a1"}, Parser: parserString, Value: reflect.ValueOf(&a1).Elem()},
			{Name: []string{"l1", "a2"}, Parser: parserString, Value: reflect.ValueOf(&a2).Elem()},
			{Name: []string{"l2", "l3", "a3"}, Parser: parserString, Value: reflect.ValueOf(&a3).Elem()},
		}},
	}

	for _, tc := range tests {
		t.Run(tc.wantDefault, func(t *testing.T) {
			a1 = "layer1"
			a2 = "layer2"
			a3 = "layer3"

			v := newFromFields(tc.inputFields, 0, `json:"%s"`)

			t.Run("Default", func(t *testing.T) {
				gotDefault, err := json.Marshal(v.Interface())
				if err != nil {
					t.Fatalf("newFromFields(%q) can't be marshalled with json: %v", tc.wantDefault, err)
				}

				if diff := cmp.Diff(string(gotDefault), tc.wantDefault); diff != "" {
					t.Errorf("newFromFields(%q) diff: (-got, +want)\n%s", tc.wantDefault, diff)
				}
			})

			t.Run("JSON", func(t *testing.T) {
				err := json.Unmarshal([]byte(tc.checkString), v.Interface())
				if err != nil {
					t.Fatalf("newFromFields(%q) can't be unmarshalled with json: %v", tc.checkString, err)
				}

				got, err := json.Marshal(v.Interface())
				if err != nil {
					t.Fatalf("newFromFields(%q) can't be marshalled with json: %v", tc.checkString, err)
				}

				if diff := cmp.Diff(string(got), tc.checkString); diff != "" {
					t.Errorf("newFromFields(%q) diff: (-got, +want)\n%s", tc.checkString, diff)
				}
			})
		})
	}
}
