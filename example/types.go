package example

//go:generate go run ./gen/...

import (
	"context"
)

type CustomType struct {
	Hell int
}
type EchoRequest struct {
	Session string `header:"session"`
	Text    string `json:"text"`
}

// EchoResponse.
type EchoResponse struct {
	Text string `json:"text"` // field comment.
}

type HelloRequest struct {
	Key string `query:"key"`
}

type HelloResponse struct {
	Session string `header:"session"`
}

type IEchoService interface {
	Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error)
	Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error)
}
