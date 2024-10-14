package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type CentralServer struct {
	inputChan chan DeliverMsg
}

func NewCentralServer() *CentralServer {
	return &CentralServer{
		inputChan: make(chan DeliverMsg),
	}
}

func (cs *CentralServer) Start() {
	fmt.Println("Starting central server...")
	for {
		deliverMsg, ok := <-cs.inputChan
		if !ok {
			continue
		}
		msg, err := json.Marshal(deliverMsg.Msg)
		if err != nil {
			continue
		}
		fmt.Println("Received message from", deliverMsg.Msg.Sender, ":", string(msg))
		if _, err = deliverMsg.Dest.Write(msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing to personal server:", err)
		}
		fmt.Println("Send Msg successfully to personal server")
	}
}
