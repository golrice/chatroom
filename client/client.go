package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	serverIp   = "172.29.218.119"
	serverPort = "8080"
)

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

	fmt.Print("Enter your message: ")

	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 1024)

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

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing message to server:", err)
			return
		}

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading response from server:", err)
			return
		}

		fmt.Println(string(buf[:n]))
	}
}
