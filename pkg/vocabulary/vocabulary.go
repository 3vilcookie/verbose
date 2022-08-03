package vocabulary

type Translation struct {
	Words []string `json:"words"`
}

type Entry map[string]Translation

type Vocabulary struct {
	Entries []Entry `json:"words"`

	Filename string
}

func New() *Vocabulary {
	v := new(Vocabulary)
	v.Entries = make(Entries)
	return v
}

func FromFile(string filename) (*Vocabulary, error) {}
func (v *Vocabulary) Load() error                   {}

func (v *Vocabulary) Save() error                      {}
func (v *Vocabulary) SaveToFile(string filename) error {}
