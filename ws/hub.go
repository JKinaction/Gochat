package ws

import (
	"encoding/json"
	"log"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte
	priMes    chan []byte
	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		priMes:     make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			count := len(h.clients)
			msg := Msg{
				Status: 5,
				Data:   map[string]interface{}{"usersCount": count},
			}
			marshal, err := json.Marshal(msg)
			if err != nil {
				log.Println(err)
				continue
			}
			h.broadcast <- marshal
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				count := len(h.clients)
				msg := Msg{
					Status: 5,
					Data:   map[string]interface{}{"usersCount": count},
				}
				marshal, err := json.Marshal(msg)
				if err != nil {
					log.Println(err)
					continue
				}
				h.broadcast <- marshal
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
		case message := <-h.priMes:
			var clientMes Msg
			json.Unmarshal(message, &clientMes)
			for client := range h.clients {
				if client.username == clientMes.Data.(map[string]interface{})["toUser"].(string) ||
					client.username == clientMes.Data.(map[string]interface{})["user"].(string) {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
