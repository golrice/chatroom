package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

// 函数表
type FunctionTable map[string]func(*PersonalServer, *Message) error

var functionTable = FunctionTable{
	"login":  Login,
	"logout": Logout,
	"send":   SendMsg,
}

// 全局服务器表
var serverTable = make(map[string]*PersonalServer)

// 全局用户密码表
var passwordTable = make(map[string]string)

// 全局中心服务器
var centerServer = NewCentralServer()

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Listen error:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is running on port 8080")

	go centerServer.Start()

	for {
		conn, err := listener.Accept()
		fmt.Println("New connection: ", conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

// login流程
func LoginFlow(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return "", errors.New("read error")
		}

		// 解析请求
		var req Request
		err = json.Unmarshal(buf[:n], &req)
		if err != nil {
			jsondata, err := genMsg("Error: parse error")

			if err != nil {
				fmt.Println("Marshal error:", err)
				continue
			}

			conn.Write(jsondata)
			continue
		}

		if req.Action != "login" {
			jsondata, err := genMsg("Error: no such action")

			if err != nil {
				fmt.Println("Marshal error:", err)
				continue
			}

			conn.Write(jsondata)
			continue
		}

		// 调用登录函数
		function := functionTable[req.Action]

		if err := function(nil, &req.Msg); err != nil {
			jsondata, err := genMsg(err.Error())

			if err != nil {
				fmt.Println("Marshal error:", err)
				continue
			}

			conn.Write(jsondata)
			continue
		}

		jsondata, err := genMsg("login Success")
		if err != nil {
			fmt.Println("Marshal error:", err)
			continue
		}

		conn.Write(jsondata)
		return req.Msg.Sender, nil
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// 消息缓冲区
	buf := make([]byte, 1024)

	// 首先处理好用户的登录，获得用户的用户名
	username, err := LoginFlow(conn)
	if err != nil {
		jsondata, err := genMsg(err.Error())

		if err != nil {
			fmt.Println("Marshal error:", err)
			return
		}

		conn.Write(jsondata)
		return
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
			fmt.Println("client is offline")
			Logout(personal_server, nil)
			break
		}

		// 解析信息
		var req Request
		err = json.Unmarshal(buf[:n], &req)
		if err != nil {
			jsondata, err := genMsg("Error: parse error")

			if err != nil {
				fmt.Println("Marshal error:", err)
				continue
			}

			conn.Write(jsondata)
			break
		}

		// 调用对应的函数
		function, ok := functionTable[req.Action]
		if !ok {
			jsondata, err := genMsg("Error: no such action")

			if err != nil {
				fmt.Println("Marshal error:", err)
				continue
			}

			conn.Write(jsondata)
			break
		}

		if err := function(personal_server, &req.Msg); err != nil {
			jsondata, err := genMsg(err.Error())

			if err != nil {
				fmt.Println("Marshal error:", err)
				continue
			}

			conn.Write(jsondata)
			continue
		}

		jsondata, err := genMsg(req.Action + " Success")
		if err != nil {
			fmt.Println("Marshal error:", err)
			continue
		}

		conn.Write(jsondata)

		// logout的时候需要重新分配server
		if req.Action == "logout" {
			username, err = LoginFlow(conn)
			if err != nil {
				jsondata, err := genMsg(err.Error())

				if err != nil {
					fmt.Println("Marshal error:", err)
					continue
				}

				conn.Write(jsondata)
				break
			}

			personal_server, ok = serverTable[username]
			if !ok {
				personal_server = NewPersonalServer(conn, make(map[string][]Message))
				serverTable[username] = personal_server
			}
		}
	}
}
