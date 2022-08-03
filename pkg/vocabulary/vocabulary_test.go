package vocabulary

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	require.NotNil(t, New())
}

func TestFromFile(t *testing.T) {
	filename := "words.json"

	tempFile, err := os.CreateTemp("", filename)
	require.Nil(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(`{
	"verbose" : { 
	  "words" : [ 
		"wortreich", 
		"langatmig", 
		"ausführlich", 
		"weitschweifig"
	  ]
	}}`))
	require.Nil(t, err)

	v, err := FromFile(tempFile.Name())
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, tempFile.Name(), v.Filename)

	words, ok := v.Entries["verbose"]
	require.True(t, ok)
	require.Equal(t, Translation{
		Words: []string{
			"wortreich",
			"langatmig",
			"ausführlich",
			"weitschweifig",
		},
	}, words)
}

func TestFromFileNotExistingFile(t *testing.T) {
	_, err := FromFile("")
	require.NotNil(t, err)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestFromFileBadJSON(t *testing.T) {
	tempFile, err := os.CreateTemp("", "bad.json")
	require.Nil(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(`{`))
	require.Nil(t, err)

	_, err = FromFile(tempFile.Name())
	require.NotNil(t, err)
	require.Error(t, err, "unexpected end of json input")
}

func TestSave(t *testing.T) {
	tempFile, err := os.CreateTemp("", "bad.json")
	require.Nil(t, err)
	defer os.Remove(tempFile.Name())
	require.Nil(t, tempFile.Close())

	v := New()
	v.Entries["assert"] = Translation{
		Words: []string{
			"behaupten",
			"versichern",
		},
	}

	require.Nil(t, v.SaveToFile(tempFile.Name()))

	content, err := os.ReadFile(tempFile.Name())
	require.Nil(t, err)
	require.Equal(
		t,
		`{"assert":{"words":["behaupten","versichern"]}}`,
		string(content),
	)
}
