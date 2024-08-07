package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterPage(server *gin.Engine) {
	server.GET("/register_page", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", nil)
	})

	server.POST("/register", func(ctx *gin.Context) {
		var credentials map[string]string
		if err := ctx.BindJSON(&credentials); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		username := credentials["username"]
		password := credentials["password"]
		email := credentials["email"]

		// query database to check if username or email already exists
		// if not, create new user and add to database
		// if yes, return error message
		if db.existsUser(username) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Username already exists"})
			return
		}

		if err := db.createUser(username, password, email); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.Redirect(http.StatusSeeOther, "/")
	})
}
