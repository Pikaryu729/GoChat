package main

import (
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"encoding/json"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request)bool {
		return true
	},
}

var serverPort = ":8080"

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func connectEndpoint(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}

	// read initial message to get name
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	var NameData struct {
		Name string `json:"name"`
	} 
	// get name from websocket message
	if err := json.Unmarshal(msg, &NameData); err != nil {
		log.Println("INVALID JSON: ", err)
		return
	}

	if NameData.Name == "" {
		log.Println("Missing Name Field")
		return
	}

	client := NewClient(NameData.Name, hub, conn)
	client.hub.register <- client
	go client.readPump()
	go client.writePump()

	log.Println("Client Successfully Connected...")
}

func setupRoutes(hub *Hub) {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		connectEndpoint(hub, w, r)
	})
}

func main() {
	hub := newHub()
	go hub.run()

	setupRoutes(hub)
	log.Printf("Server is listening on port %s\n", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}