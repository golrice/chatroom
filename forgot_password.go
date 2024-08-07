package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func forgotPasswordPage(server *gin.Engine) {
	server.GET("/forgot-password_page", func(ctx *gin.Context) {
		ctx.HTML(http.StatusMovedPermanently, "forgot-password.html", nil)
	})

	server.POST("/forgot-password", func(ctx *gin.Context) {
		// get email from frontend
		var credentials map[string]string
		if err := ctx.BindJSON(&credentials); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		email := credentials["email"]

		// send email to user
		// get return state from user
		if ok := sendEmail(email); ok {
			// if success, redirect to register
			ctx.Redirect(http.StatusSeeOther, "/register_page")
			return
		}
		// else send a json
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "fail to send email",
		})
	})
}

func sendEmail(email string) bool {
	return len(email) != 0
}
