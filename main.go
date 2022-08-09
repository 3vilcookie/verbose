package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/RaphaelPour/verbose/pkg/vocabulary"
	"github.com/gin-gonic/gin"
)

var (
	//go:embed index.tmpl
	indexTemplate string

	//go:embed logo.png
	logoImage []byte

	voc vocabulary.Vocabulary

	Filename = flag.String("vocabulary-file", "vocabulary.json", "Path to vocabulary file.")
	Port     = flag.Int("port", 8080, "Serverport")
)

func main() {
	flag.Parse()

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
	router.SetHTMLTemplate(renderer)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", voc.Entries)
	})

	router.GET("logo.png", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/png", logoImage)
	})

	router.POST("/", func(c *gin.Context) {
		en, ok := c.GetPostForm("en")
		if !ok {
			c.String(http.StatusBadRequest, "parameter 'en' missing")
			return
		}

		rawDeList, ok := c.GetPostForm("de")
		if !ok {
			c.String(http.StatusBadRequest, "parameter 'de' missing")
			return
		}

		deList := strings.Split(rawDeList, ", ")

		if _, exists := voc.Entries[en]; exists {
			c.String(http.StatusBadRequest, "duplicate, english word already existing")
			return
		}

		voc.Entries[en] = vocabulary.Translation{
			Words: deList,
		}

		voc.Save()

		// redirect to index
		c.Request.URL.Path = "/"
		c.Request.Method = "GET"
		router.HandleContext(c)
	})

	if err := router.Run(fmt.Sprintf("localhost:%d", *Port)); err != nil {
		fmt.Println(err)
	}
}
