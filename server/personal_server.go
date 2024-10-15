package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
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
	// 本地持久化存储消息
	p.MsgStore[msg.Receiver] = append(p.MsgStore[msg.Receiver], *msg)

	// 检查当前Recevier是否是群聊，如果是群聊需要将消息转发到群聊中
	if _, exists := groupTable[msg.Receiver]; exists {
		// 群聊转发消息
		for _, user := range groupTable[msg.Receiver] {
			// 向其他用户转发消息
			deliverMsg := DeliverMsg{
				Msg:  *msg,
				Dest: serverTable[user].Conn,
			}

			centerServer.inputChan <- deliverMsg
		}
	} else {
		// 查看当前的Recevier是否存在
		if _, exists := serverTable[msg.Receiver]; !exists {
			return errors.New(msg.Receiver + " not exist")
		}
		serverTable[msg.Receiver].MsgStore[msg.Sender] = append(serverTable[msg.Receiver].MsgStore[msg.Sender], *msg)
		// 向其他用户转发消息
		deliverMsg := DeliverMsg{
			Msg:  *msg,
			Dest: serverTable[msg.Receiver].Conn,
		}

		centerServer.inputChan <- deliverMsg
	}

	return nil
}

func ShowOther(p *PersonalServer, msg *Message) error {
	// 展示目前有什么用户在登录
	otherUsers := make([]string, 0)

	if msg.Receiver == "all" {
		for key, value := range serverTable {
			if value.Conn != nil {
				otherUsers = append(otherUsers, key)
			}
		}
	} else {
		// 说明是想查找群聊中的成员
		if _, exists := groupTable[msg.Receiver]; exists {
			otherUsers = groupTable[msg.Receiver]
		}
	}

	retMessgae := Message{
		Sender:   "server",
		Receiver: msg.Sender,
		Content:  strings.Join(otherUsers, ","),
	}
	deliverMsg := DeliverMsg{
		Msg:  retMessgae,
		Dest: p.Conn,
	}

	centerServer.inputChan <- deliverMsg

	return nil
}

func JoinGroup(p *PersonalServer, msg *Message) error {
	// 检查是否存在当前群聊，如果不存在需要make
	if _, exists := groupTable[msg.Receiver]; !exists {
		groupTable[msg.Receiver] = make([]string, 0)
	}

	groupTable[msg.Receiver] = append(groupTable[msg.Receiver], msg.Sender)

	return nil
}

func DisplayHistory(p *PersonalServer, msg *Message) error {
	// 展示历史消息
	if _, exists := p.MsgStore[msg.Receiver]; !exists {
		return errors.New("no history message")
	}

	historyMsg := p.MsgStore[msg.Receiver]
	if len(historyMsg) == 0 {
		return errors.New("no history message")
	}

	historyMsgStr := make([]string, 0)
	for _, msg := range historyMsg {
		historyMsgStr = append(historyMsgStr, msg.Content)
	}

	retMessgae := Message{
		Sender:   "server",
		Receiver: msg.Sender,
		Content:  strings.Join(historyMsgStr, "\n"),
	}
	deliverMsg := DeliverMsg{
		Msg:  retMessgae,
		Dest: p.Conn,
	}

	centerServer.inputChan <- deliverMsg

	return nil
}
