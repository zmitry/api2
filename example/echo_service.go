package example

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type EchoService struct {
	sessions map[string]bool
}

func NewEchoService() *EchoService {
	return &EchoService{
		sessions: make(map[string]bool),
	}
}

func (s *EchoService) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	if req.Key != "secret password" {
		return nil, fmt.Errorf("bad key")
	}
	sessionBytes := make([]byte, 16)
	if _, err := rand.Read(sessionBytes); err != nil {
		return nil, err
	}
	session := hex.EncodeToString(sessionBytes)
	s.sessions[session] = true
	return &HelloResponse{
		Session: session,
	}, nil
}

func (s *EchoService) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	if !s.sessions[req.Session] {
		return nil, fmt.Errorf("bad session")
	}
	return &EchoResponse{
		Text: req.Text,
	}, nil
}
