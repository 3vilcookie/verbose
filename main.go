package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/RaphaelPour/verbose/pkg/vocabulary"
	"github.com/gin-gonic/gin"
)

var (
	//go:embed index.tmpl
	indexTemplate string

	//go:embed logo.png
	logoImage []byte

	voc vocabulary.Vocabulary

	mutex sync.Mutex

	Filename        = flag.String("vocabulary-file", "vocabulary.json", "Path to vocabulary file.")
	Port            = flag.Int("port", 8080, "Serverport")
	CredentialsFile = flag.String("credentials-file", "credentials.json", "Path to credentials file with user and pw.")
)

func main() {
	flag.Parse()

	// parse credentials file
	accounts, err := loadAccounts(*CredentialsFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// add pipeline function 'join' to convert en-words to list
	renderFunctions := template.FuncMap{
		"join": strings.Join,
	}

	// set up custom renderer in order to use an embedded template
	renderer, err := template.New("index").Funcs(renderFunctions).Parse(indexTemplate)
	if err != nil {
		fmt.Println(err)
		return
	}

	var voc *vocabulary.Vocabulary
	if _, err := os.Stat(*Filename); os.IsNotExist(err) {
		// create new vocabulary if file not existing
		voc = vocabulary.New()
		voc.Filename = *Filename
		voc.Save()
	} else {
		// load vocabulary
		if voc, err = vocabulary.LoadFile(*Filename); err != nil {
			fmt.Printf("error loading vocabulary from file '%s': %s\n", *Filename, err)
			return
		}
	}

	router := gin.Default()

	// initialize basic auth for /new routes
	authorized := router.Group("/", gin.BasicAuth(accounts))

	router.SetHTMLTemplate(renderer)
	router.GET("/", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		c.HTML(http.StatusOK, "index", voc.Entries)
	})

	router.GET("logo.png", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/png", logoImage)
	})

	authorized.POST("/new", func(c *gin.Context) {
		en, ok := c.GetPostForm("en")
		if !ok || en == "" {
			c.String(http.StatusBadRequest, "parameter 'en' missing")
			return
		}

		rawDeList, ok := c.GetPostForm("de")
		if !ok || rawDeList == "" {
			c.String(http.StatusBadRequest, "parameter 'de' missing")
			return
		}

		deList := strings.Split(rawDeList, ",")

		if _, exists := voc.Entries[en]; exists {
			c.String(http.StatusBadRequest, "duplicate, english word already existing")
			return
		}

		mutex.Lock()
		defer mutex.Unlock()
		voc.Entries[en] = vocabulary.Translation{
			Words: deList,
		}
		voc.Save()

		c.Redirect(http.StatusFound, "/")
	})

	if err := router.Run(fmt.Sprintf("localhost:%d", *Port)); err != nil {
		fmt.Println(err)
	}
}

func loadAccounts(filename string) (gin.Accounts, error) {
	if filename == "" {
		return nil, errors.New("cannot start without users file")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", filename, err)
	}

	users := make(gin.Accounts)
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}

	return users, nil
}
