package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Listen error:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is running on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		fmt.Println("Received message:", string(buf[:n]))

		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}
