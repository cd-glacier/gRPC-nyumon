package main

import (
	"context"

	"github.com/g-hyoga/gRPC-nyumon/echo"
)

type echoService struct{}

func (s *echoService) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Message: req.GetMessage()}, nil
}
