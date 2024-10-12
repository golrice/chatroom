package main

import (
	"errors"
	"net"
	"strings"
)

// 负责消息的存储以及转发
type PersonalServer struct {
	// id
	conn net.Conn
	// 消息存储器
	msgStore MessageStore
}

type Message struct {
	sender  string `json:"sender"`
	content string `json:"content"`
}

// 消息解析函数
func ParseMsg(msg string) (string, []string) {
	// action args...
	s := strings.Split(msg, " ")
	return s[0], s[1:]
}

// 历史信息队列
type HistoryMsgQueue = []Message

// 消息存储器
type MessageStore = map[string]HistoryMsgQueue

func NewPersonalServer(conn net.Conn, msgStore MessageStore) *PersonalServer {
	return &PersonalServer{
		conn:     conn,
		msgStore: msgStore,
	}
}

// 消息处理
func Login(p *PersonalServer, args []string) error {
	// 查看当前是否已经登录
	if p != nil {
		return errors.New("already login")
	}

	// 解析登录信息
	username, password := args[0], args[1]

	// 验证登录信息
	if serverTable[username].conn != nil {
		return errors.New("username already exist")
	}

	pwd, ok := passwordTable[username]
	if !ok {
		// 说明是第一次登录
		passwordTable[username] = password
		return nil
	}

	if pwd != password {
		return errors.New("password error")
	}

	// 登录成功
	return nil
}

func Logout(p *PersonalServer, args []string) error {
	p.conn = nil

	return nil
}

func SendMsg(p *PersonalServer, args []string) error {
	return nil
}
