package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	serverIp   = "51.13.118.186"
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

func printWelcome() {
	color.Cyan("Welcome to the chat client!")
	color.Cyan("Type 'login <username> <password>' to login.")
	color.Cyan("Type 'send <receiver> <message>' to send a message.")
	color.Cyan("Type 'show [all | username]' to display messages.")
	color.Cyan("Type 'group <groupname>' to interact with a group.")
	color.Cyan("Type 'display <option>' to control display options.")
	color.Cyan("Type 'exit' to exit the chat.\n")
}

func printPrompt(username string) {
	time.Sleep(time.Second)
	color.Green("[%s] >> ", username)
}

func formatMessage(sender, content string) {
	color.Blue("%s [%s]:\n", sender, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(content)
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
		color.Red("Error connecting to server: %v", err)
		return
	}
	defer conn.Close()

	printWelcome()

	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 1024)

	username := "client"

	go func() {
		for {
			n, err := conn.Read(buf)
			if err != nil {
				color.Red("Server disconnected: %v", err)
				os.Exit(1)
			}

			var response Message
			if err := json.Unmarshal(buf[:n], &response); err != nil {
				color.Red("Error decoding JSON: %v", err)
			}
			formatMessage(response.Sender, fmt.Sprintf("%s:%s", response.Sender, response.Content))
		}
	}()

	for {
		printPrompt(username)
		message, err := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if err != nil {
			color.Red("Error reading message: %v", err)
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

		switch action {
		case "login":
			if len(args) < 3 {
				color.Red("Usage: login <username> <password>")
				continue
			}
			username = args[1]
			msg.Sender = username
			msg.Content = args[2]
			req.Msg = msg
		case "send":
			if len(args) < 3 {
				color.Red("Usage: send <receiver> <message>")
				continue
			}
			msg.Receiver = args[1]
			msg.Content = strings.Join(args[2:], " ")
			req.Msg = msg
		case "show":
			msg.Receiver = "all"
			if len(args) == 2 {
				msg.Receiver = args[1]
			}
			req.Msg = msg
		case "group", "display":
			if len(args) != 2 {
				color.Red("Invalid %s command", action)
				continue
			}
			msg.Receiver = args[1]
			req.Msg = msg
		case "logout":
			color.Yellow("Logging out...")
		default:
			color.Red("Invalid action: %s", action)
			continue
		}

		jsondata, err := json.Marshal(req)
		if err != nil {
			color.Red("Error encoding JSON: %v", err)
			return
		}

		_, err = conn.Write(jsondata)
		if err != nil {
			color.Red("Error writing to server: %v", err)
			return
		}
	}
}
