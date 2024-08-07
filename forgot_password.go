package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func forgotPasswordPage(server *gin.Engine) {
	server.GET("/forgot-password", func(ctx *gin.Context) {
		ctx.HTML(http.StatusMovedPermanently, "forgot-password.html", nil)
	})
}
