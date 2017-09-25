//go:generate gooption-cli -p service -r Price -r Greek -r ImpliedVol
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func newEngine() *gin.Engine {
	// Starts a new Gin instance with no middle-ware
	router := gin.Default()
	router.POST("/price", handlerPrice)
	router.POST("/greek", handlerGreek)
	router.POST("/impliedvol", handlerImpliedVol)
	return router
}

func main() {
	newEngine().Run()
}

// This function's name is a must. App Engine uses it to drive the requests properly.
func init() {
	// Handle all requests using net/http
	http.Handle("/", newEngine())
}
