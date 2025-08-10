//go:build ignore
// +build ignore

package main

import (
	"log"
	"net/http"
	"server/server"
)

func main() {
	ws_server := server.NewServer()
	http.HandleFunc("/ws", ws_server.HandleConnections)
	log.Println("http server started on :8008")
	err := http.ListenAndServe(":8008", nil)
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
}
