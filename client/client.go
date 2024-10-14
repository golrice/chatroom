package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	serverIp = "172.29.218.119"
	// serverIp   = "51.13.118.186"
	serverPort = "8080"
)

type Message struct {
	Sender   string `json:"Sender"`
	Receiver string `json:"Receiver"`
	Content  string `json:"Content"`
}

type Request struct {
	Action string  `json:"Action"`
	Msg    Message `json:"Message"`
}

func main() {
	if len(os.Args) > 1 {
		serverIp = os.Args[1]
	}
	if len(os.Args) > 2 {
		serverPort = os.Args[2]
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", serverIp, serverPort))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Enter your message:")

	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 1024)

	username := "client"

	go func() {
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("server disconnected: ", err)
				os.Exit(1)
			}

			var response Message
			if err := json.Unmarshal(buf[:n], &response); err != nil {
				fmt.Println("Error decoding JSON:", err)
			}
			fmt.Println(response.Content)
		}
	}()

	for {
		message, err := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		if message == "exit" {
			break
		}

		args := strings.Split(message, " ")
		action := args[0]

		msg := Message{
			Sender:   username,
			Receiver: "server",
			Content:  "",
		}
		req := Request{
			Action: action,
			Msg:    msg,
		}

		if action == "login" {
			username = args[1]
			msg.Sender = username
			msg.Content = args[2]
			req.Msg = msg
		} else if action == "send" {
			msg.Receiver = args[1]
			msg.Content = strings.Join(args[2:], " ")
			req.Msg = msg
		} else if action != "logout" {
			fmt.Println("Invalid action:", action)
			continue
		}

		jsondata, err := json.Marshal(req)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		_, err = conn.Write(jsondata)
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}
	}
}
