package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// create a default server
	ginServer := NewServer("127.0.0.1", 8080)
	var err error
	db, err = NewDb("127.0.0.1", 3306, "chatroom")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Database.Close()

	// init files path
	ginServer.engine.Static("static", "./static")
	ginServer.engine.LoadHTMLGlob("./templates/*")

	EmailVerification(ginServer.engine)

	LoginPage(ginServer.engine)
	RegisterPage(ginServer.engine)
	forgotPasswordPage(ginServer.engine)

	ChatPage(ginServer.engine)

	ginServer.engine.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "404.html", nil)
	})

	// apply a port for server
	ginServer.Run()
}
