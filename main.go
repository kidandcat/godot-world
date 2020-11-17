package main

import (
	"github.com/gin-gonic/gin"
)

var db *DB

func main() {
	db = dbInitialize()
	r := gin.Default()
	m := serverSetupRoutes()

	r.GET("/", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	r.Run(":5000")
}
