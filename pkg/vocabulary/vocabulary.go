package vocabulary

import (
	"encoding/json"
	"os"
)

type Example struct {
	Original    string `json:"original,omitempty"`
	Translation string `json:"translation,omitempty"`
}

type Translation struct {
	Words   []string `json:"words"`
	Example *Example `json:"example,omitempty"`
}

type Vocabulary struct {
	Entries map[string]Translation

	Filename string
}

func New() *Vocabulary {
	v := new(Vocabulary)
	v.Entries = make(map[string]Translation)
	return v
}

func LoadFile(filename string) (*Vocabulary, error) {
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

func (v *Vocabulary) SaveFile(filename string) error {
	v.Filename = filename
	return v.Save()
}
