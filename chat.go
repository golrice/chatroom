package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChatPage(server *gin.Engine) {
	server.GET("/chat", func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			username = "Anonymous"
		}
		c.HTML(http.StatusOK, "chat.html", gin.H{
			"username": username,
		})
	})
}
