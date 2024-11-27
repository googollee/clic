package sources

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googollee/go-cfg/structtags"
)

func TestFindEndFieldHaveSameIndex(t *testing.T) {
	tests := []struct {
		inputNames [][]string
		inputIndex int
		want       int
	}{
		{[][]string{}, 0, 0},
		{[][]string{{}}, 0, 0},

		{[][]string{{"a1"}, {"a2"}, {"a3"}}, 0, 3},
		{[][]string{{"a1"}, {"a2"}, {"a3"}}, 1, 0},

		{[][]string{{"a1", "b1"}, {"a1", "b2"}, {"a2"}}, 0, 3},
		{[][]string{{"a1", "b1"}, {"a1", "b2"}, {"a2"}}, 1, 2},
		{[][]string{{"a1", "b1"}, {"a1", "b2"}, {"a2", "b3"}}, 1, 3},
		{[][]string{{"a1", "b1"}, {"a1", "b2", "c1"}, {"a2", "b3"}}, 1, 3},

		{[][]string{{"a1", "b1"}, {"a2"}, {"a3", "b3"}}, 1, 1},
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

func TestNewFromFields(t *testing.T) {
	tests := []struct {
		wantDefault string
		checkString string
		inputFields []structtags.Field
	}{
		{`{"a1":"abc"}`, `{"a1":"123"}`, []structtags.Field{
			{Name: []string{"a1"}, Value: reflect.ValueOf("abc")},
		}},
	}

	for _, tc := range tests {
		t.Run(tc.wantDefault, func(t *testing.T) {
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
