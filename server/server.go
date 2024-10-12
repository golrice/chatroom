package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// 函数表
type FunctionTable map[string]func(*PersonalServer, []string) error

var functionTable = FunctionTable{
	"login":  Login,
	"logout": Logout,
	"send":   SendMsg,
}

// 全局服务器表
var serverTable = make(map[string]*PersonalServer)

// 全局用户密码表
var passwordTable = make(map[string]string)

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
	// 消息缓冲区
	buf := make([]byte, 1024)

	// 首先处理好用户的登录，获得用户的用户名
	username := ""
	for {
		n, err := conn.Read(buf)
		if err != nil {
			conn.Write([]byte("Error: read error"))
			return
		}

		// 解析信息
		var msg Message
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			conn.Write([]byte("Error: parse error"))
			continue
		}

		// 调用登录函数
		action, args := ParseMsg(msg.content)
		if action != "login" {
			conn.Write([]byte("Error: not login"))
			continue
		}

		username = args[0]
		function := functionTable[action]
		if err := function(nil, args); err != nil {
			conn.Write([]byte("Error: " + err.Error()))
			continue
		}
		break
	}

	// 生成一个用户端处的 server 负责存储信息以及信息的转发
	personal_server, ok := serverTable[username]
	if !ok {
		personal_server = NewPersonalServer(conn, make(map[string][]Message))
		serverTable[username] = personal_server
	}

	for {
		// 读取信息
		n, err := conn.Read(buf)
		if err != nil {
			Logout(personal_server, nil)
			break
		}

		// 解析信息
		var msg Message
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			conn.Write([]byte("Error: parse error"))
			break
		}

		// 调用对应的函数
		action, args := ParseMsg(msg.content)
		function := functionTable[action]

		if function == nil {
			conn.Write([]byte("Error: no such action"))
			continue
		}

		if err := function(personal_server, args); err != nil {
			conn.Write([]byte("Error: " + err.Error()))
			continue
		}
	}
}
