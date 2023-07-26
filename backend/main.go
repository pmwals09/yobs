package backend

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()
	router.GET("/", getHome)
	router.Run("localhost:8080")
}

func getHome(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello")
}
