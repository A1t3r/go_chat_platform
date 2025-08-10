package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader *websocket.Upgrader
	clients  []*websocket.Conn
}

func NewServer() *Server {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
	}
	var clients []*websocket.Conn
	return &Server{&upgrader, clients}
}

func (server *Server) ping() {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()
	for range ticker.C {
		var updated_clients []*websocket.Conn
		for _, c := range server.clients {
			log.Println("Ping", c.RemoteAddr())
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				log.Println("Deleting", c.RemoteAddr())
			} else {
				c.SetWriteDeadline(time.Now().Add(20 * time.Second))
				//TODO:  swap element to the end and delete tail
				updated_clients = append(updated_clients, c)
			}
		}
		server.clients = updated_clients
	}
}

func (server *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	server.clients = append(server.clients, ws)
	defer ws.Close()

	ws.SetPongHandler(func(string) error {
		log.Println("Got pong from", ws.RemoteAddr())
		ws.SetReadDeadline(time.Now().Add(20 * time.Second))
		return nil
	})
	go server.ping()

	// main loop
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(string(message))
		ws.SetReadDeadline(time.Now().Add(90 * time.Second))

		// should be probably goroutine as well
		actual_clients := append([]*websocket.Conn{}, server.clients...)
		for _, client := range actual_clients {
			if err := client.WriteMessage(messageType, message); err != nil {
				log.Println(err)
			}
		}
	}
}
