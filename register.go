package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterPage(server *gin.Engine) {
	server.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusMovedPermanently, "register.html", nil)
	})
}
