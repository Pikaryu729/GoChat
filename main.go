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

func connectEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {return true}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Successfully Connected...")

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(msg))

	var NameData struct {
		Name string `json:"name"`
	} 

	if err := json.Unmarshal(msg, &NameData); err != nil {
		log.Println("INVALID JSON: ", err)
		return
	}

	if NameData.Name == "" {
		log.Println("Missing Name Field")
		return
	}
	client := NewClient(NameData.Name, conn)
	log.Printf("User Connected: %s", client.Name)
	reader(conn)
}



func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/connect", connectEndpoint)
}

func main() {
	setupRoutes()
	log.Printf("Server is listening on port %s\n", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}