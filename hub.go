package main

import (
	"fmt"
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	host       chan *Message
}

type Message struct {
	BoardId int
	Body    string
}

func newMessage(id int, body string) *Message {
	return &Message{
		BoardId: id,
		Body:    body,
	}
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		host:       make(chan *Message),
		clients:    make(map[*Client]bool),
	}
}

// Send a message to boardId's host
func (h *Hub) sendMsg(msg *Message) {
	for client := range h.clients {
		if client.BoardId == msg.BoardId {
			log.Println(fmt.Sprintf("sending message to client %d", client.BoardId))
			client.send <- []byte(msg.Body)
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case msg := <-h.host:
			h.sendMsg(msg)
		}
	}
}
