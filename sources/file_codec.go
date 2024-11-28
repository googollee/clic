package sources

import (
	"encoding/json"
	"os"
)

type JSON struct{}

func (JSON) TagName() string {
	return "json"
}

func (JSON) ExtName() string {
	return "json"
}

func (JSON) Encode(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(v)
}

func (JSON) Decode(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(v)
}
