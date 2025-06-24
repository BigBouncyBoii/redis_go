package main

import (
	"github.com/BigBouncyBoii/redis_go/internal/server"
	"github.com/BigBouncyBoii/redis_go/internal/protocol"
	"fmt"
)

func main() {
	s := server.NewServer(":6969")
	db := server.NewDatabase()
	go func(){
		for msg := range s.GetMsgChannel() {
			go func(msg server.Message) {
				fmt.Printf("Received message: %v from %v", string(msg.Payload), msg.From)
				args, err := protocol.ParseCommand(msg.Payload)
				if err != nil{
					fmt.Printf("Error parsing command: %v", err)
					s.WriteError(msg.Conn, err)
					return
				}
				response, err := server.HandleCommand(args, db)
				if err != nil {
					fmt.Printf("Error handling command: %v", err)
					s.WriteError(msg.Conn, err)
					return
				}
				byteResponse := protocol.ParseRespone(response)
				s.WriteResponse(msg.Conn, byteResponse)
			} (msg)
		}
	}()
	s.Start()

}