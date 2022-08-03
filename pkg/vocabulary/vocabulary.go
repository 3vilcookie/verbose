package vocabulary

import (
	"encoding/json"
	"os"
)

type Translation struct {
	Words []string `json:"words"`
}

type Vocabulary struct {
	Entries map[string]Translation `json:"words"`

	Filename string
}

func New() *Vocabulary {
	v := new(Vocabulary)
	v.Entries = make(map[string]Translation)
	return v
}

func FromFile(filename string) (*Vocabulary, error) {
	v := New()
	v.Filename = filename
	if err := v.Load(); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Vocabulary) Load() error {
	// load json from file
	rawJSON, err := os.ReadFile(v.Filename)
	if err != nil {
		return err
	}

	// unmarshal json to struct
	return json.Unmarshal(rawJSON, &v.Entries)
}

func (v *Vocabulary) Save() error {
	// marshal struct to json
	rawJSON, err := json.Marshal(v.Entries)
	if err != nil {
		return err
	}

	// write json to file
	return os.WriteFile(v.Filename, rawJSON, 0600)
}

func (v *Vocabulary) SaveToFile(filename string) error {
	v.Filename = filename
	return v.Save()
}
