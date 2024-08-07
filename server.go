package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	Ip   string
	Port int

	engine *gin.Engine

	clients   map[*websocket.Conn]bool
	broadCast chan Message
	upgrader  websocket.Upgrader
}

type Message struct {
	UserName string `json:"username"`
	Message  string `json:"message"`
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,

		engine: gin.Default(),

		clients:   make(map[*websocket.Conn]bool),
		broadCast: make(chan Message),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (server *Server) Run() {
	go server.broadCastLoop()

	server.engine.GET("/ws/chat/:username", func(ctx *gin.Context) {
		server.handleConnection(ctx.Writer, ctx.Request)
	})

	server.engine.Run(fmt.Sprintf("%s:%d", server.Ip, server.Port))
}

func (server *Server) broadCastLoop() {
	for {
		message := <-server.broadCast
		for client := range server.clients {
			if err := client.WriteJSON(message); err != nil {
				client.Close()
				delete(server.clients, client)
			}
		}
	}
}

func (server *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer ws.Close()

	server.clients[ws] = true

	for {
		var message Message
		if err := ws.ReadJSON(&message); err != nil {
			delete(server.clients, ws)
			break
		}
		server.broadCast <- message
	}
}
