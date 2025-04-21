package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID uuid.UUID
	Name string `json:"name"`
	Conn *websocket.Conn
	Send chan []byte
}

func NewClient(name string, conn *websocket.Conn) *Client{
	return &Client{
		ID: uuid.New(),
		Name: name,
		Conn: conn,
		Send: make(chan []byte),
	}
}

func (c *Client) SendMessage(msg []byte) {
	c.Send <- msg
}