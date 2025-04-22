package main

type Message struct {
	Sender string `json:"sender"`
	Body []byte `json:"body"`
}