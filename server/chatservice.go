package main

import (
	"io"
	"log"
	"sync"

	"github.com/g-hyoga/gRPC-nyumon/chat"
	"github.com/k0kubun/pp"
)

var streams sync.Map

type chatService struct{}

func (chatService) Connect(stream chat.ChatService_ConnectServer) error {
	log.Println("connect", &stream)

	streams.Store(stream, struct{}{})
	defer func() {
		log.Println("disconnect", &stream)
		streams.Delete(stream)
	}()

	for {
		req, err := stream.Recv()
		pp.Println(err)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		streams.Range(func(key, value interface{}) bool {
			stream := key.(chat.ChatService_ConnectServer)
			stream.Send(&chat.Post{
				Name:    req.GetName(),
				Message: req.GetMessage(),
			})
			return true
		})
	}
}
