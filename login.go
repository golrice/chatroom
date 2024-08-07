package main

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func LoginPage(server *gin.Engine) {
	server.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	loginAccess(server)
}

func authenticateUser(username, password string) bool {
	return db.authenticate(username, password)
}

func loginAccess(server *gin.Engine) {
	server.POST("/login", func(ctx *gin.Context) {
		var credentials map[string]string
		if err := ctx.BindJSON(&credentials); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
			return
		}
		username := credentials["username"]
		password := credentials["password"]

		// 验证用户信息
		if authenticateUser(username, password) {
			// 验证成功，重定向到聊天页面
			redirectURL := "/chat?username=" + url.QueryEscape(username)
			ctx.Redirect(http.StatusSeeOther, redirectURL)
		} else {
			// 验证失败，返回错误信息
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid username or password",
			})
		}
	})
}
