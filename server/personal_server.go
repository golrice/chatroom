package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

// 负责消息的存储以及转发
type PersonalServer struct {
	// id
	Conn net.Conn
	// 持久化消息存储器
	MsgStore MessageStore
}

type Message struct {
	Sender   string `json:"Sender"`
	Receiver string `json:"Receiver"`
	Content  string `json:"Content"`
}

type DeliverMsg struct {
	Msg  Message
	Dest net.Conn
}

// 发送服务器请求
type Request struct {
	Action string  `json:"Action"`
	Msg    Message `json:"Message"`
}

// 历史信息队列
type HistoryMsgQueue = []Message

// 消息存储器
type MessageStore = map[string]HistoryMsgQueue

func NewPersonalServer(conn net.Conn, msgStore MessageStore) *PersonalServer {
	return &PersonalServer{
		Conn:     conn,
		MsgStore: msgStore,
	}
}

func genMsg(content string) ([]byte, error) {
	msg := Message{
		Sender:   "server",
		Receiver: "client",
		Content:  content,
	}

	jsondata, err := json.Marshal(msg)
	return jsondata, err
}

// 消息处理
func Login(p *PersonalServer, msg *Message) error {
	// 解析登录信息
	username, password := msg.Sender, msg.Content

	fmt.Println(username, " logining...")

	// 验证登录信息
	if value, exists := serverTable[username]; exists && value.Conn != nil {
		return errors.New("username already exist")
	}

	pwd, ok := passwordTable[username]
	if !ok {
		// 说明是第一次登录
		passwordTable[username] = password
	} else if pwd != password {
		return errors.New("password error")
	}

	return nil
}

func Logout(p *PersonalServer, msg *Message) error {
	fmt.Println(p.Conn.RemoteAddr().String(), " logout")
	// 个人服务器下线
	p.Conn = nil

	return nil
}

func SendMsg(p *PersonalServer, msg *Message) error {
	fmt.Println(msg.Sender, " send message to ", msg.Receiver, " : ", msg.Content)
	// 本地持久化存储消息
	p.MsgStore[msg.Receiver] = append(p.MsgStore[msg.Receiver], *msg)

	// 向其他用户转发消息
	deliverMsg := DeliverMsg{
		Msg:  *msg,
		Dest: serverTable[msg.Receiver].Conn,
	}
	fmt.Println("Recevier: ", msg.Receiver)
	fmt.Println("Conn: ", serverTable[msg.Receiver].Conn.RemoteAddr().String())

	centerServer.inputChan <- deliverMsg

	return nil
}
