package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"math/rand"
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
		"random_word": func(words map[string]vocabulary.Translation) string {
			if len(words) == 0 {
				return "NULL - no words available"
			}
			random_index := rand.Intn(len(words))
			for en, translation := range words {
				if random_index > 0 {
					random_index--
					continue
				}
				return fmt.Sprintf("%s - %s", en, strings.Join(translation.Words, ","))
			}
			return ""
		},
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
		if err := voc.Save(); err != nil {
			fmt.Println(err)
			return
		}
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
		entry := vocabulary.Translation{
			Words: deList,
		}
		original, ok1 := c.GetPostForm("example_original")
		translation, ok2 := c.GetPostForm("example_translation")
		fmt.Println(original, ok1)
		fmt.Println(translation, ok2)
		if ok1 && ok2 {
			entry.Example = &vocabulary.Example{
				Original:    original,
				Translation: translation,
			}
		}

		voc.Entries[en] = entry
		if err := voc.Save(); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error saving new entry: %s", err))
			return
		}

		c.Redirect(http.StatusFound, "/")
	})

	apiv1 := router.Group("api/v1")
	apiv1.Use(func(c *gin.Context) {
		/* add cors header */
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH,OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			return
		}

		c.Next()
	})
	authorizedAPIv1 := apiv1.Group("/", gin.BasicAuth(accounts))

	apiv1.GET("/words", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		c.JSON(http.StatusOK, voc.Entries)
	})

	apiv1.GET("/words/:word", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		translation, exists := voc.Entries[c.Param("word")]
		if !exists {
			c.JSON(
				http.StatusNotFound,
				map[string]string{
					"error": "word not found",
				},
			)
			return
		}

		c.JSON(http.StatusOK, translation)
	})

	authorizedAPIv1.POST("/words/:word", func(c *gin.Context) {
		if _, exists := voc.Entries[c.Param("word")]; exists {
			c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "duplicate, english word already existing",
				},
			)
			return
		}
		var data vocabulary.Translation
		if c.BindJSON(&data) != nil {
			c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "error parsing json",
				},
			)
		}

		mutex.Lock()
		defer mutex.Unlock()
		voc.Entries[c.Param("word")] = data
		if err := voc.Save(); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error saving new entry: %s", err))
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	if err := router.Run(fmt.Sprintf(":%d", *Port)); err != nil {
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
