package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Message struct {
	sender *Client
	body []byte
}

type Client struct {
	id uuid.UUID 
	name string 
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}

func NewClient(name string, hub *Hub, conn *websocket.Conn) *Client{
	return &Client{
		id: uuid.New(),
		name: name,
		hub: hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Error Reading Message from client", err)
			break
		}
		log.Printf("Message from %s: %s", c.name, msg)
		c.hub.broadcast <- msg 
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error Writing Message")
			break
		}
	}
}