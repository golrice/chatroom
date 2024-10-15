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
	for {
		deliverMsg, ok := <-cs.inputChan
		if !ok || deliverMsg.Dest == nil {
			continue
		}
		msg, err := json.Marshal(deliverMsg.Msg)
		if err != nil {
			continue
		}
		if _, err = deliverMsg.Dest.Write(msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing to personal server:", err)
		}
	}
}
