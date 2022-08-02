package main // import "github.com/RaphaelPour/verbose"

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	})
	router.Run()
}
