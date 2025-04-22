package main

import (
	"log"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)


type Client struct {
	id uuid.UUID 
	name string 
	hub *Hub
	conn *websocket.Conn
	send chan *Message
}

func NewClient(name string, hub *Hub, conn *websocket.Conn) *Client{
	return &Client{
		id: uuid.New(),
		name: name,
		hub: hub,
		conn: conn,
		send: make(chan *Message),
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
		message := &Message{Sender: c.name, Body: msg}
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for message := range c.send {
		jsonData, err := json.Marshal(message)
		if err != nil {
			log.Println("error marshaling message into json")
			break
		}
		err = c.conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("Error Writing Message")
			break
		}
	}
}