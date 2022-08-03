package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
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
)

func main() {
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

	// load vocabulary
	voc := vocabulary.New()
	voc.Entries["verbose"] = vocabulary.Translation{
		Words: []string{
			"wortreich",
			"langatmig",
			"ausführlich",
			"weitschweifig",
		},
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

		// redirect to index
		c.Request.URL.Path = "/"
		c.Request.Method = "GET"
		router.HandleContext(c)
	})

	if err := router.Run(); err != nil {
		fmt.Println(err)
	}
}
