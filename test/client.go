package main

import (
	"github.com/funny/link/codec"
	"github.com/funny/link"
	"fmt"
)


type PtpvMessage struct {
	Cmd string `json:"cmd"`
	From string `json:"from"`
	To   string `json:"to"`
	Channel string `json:"channel"`
	Body string `json:"body"`
}

func main()  {
	json := codec.Json()
	json.RegisterName("PtpvMessage",PtpvMessage{})

	client, err := link.Dial("tcp", "localhost:8787", json, 0)
	if err != nil {
		fmt.Println("dial error")
	}
	msg := PtpvMessage{Cmd:"REGISTER",From:"1003",To:"10004",Body:"1234"}
	client.Send(msg)

}