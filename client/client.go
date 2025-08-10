package client

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	host        string
	name        string
	time_format string
	ws          *websocket.Conn
}

func NewClient(host string, name string) *Client {
	u := url.URL{Scheme: "ws", Host: host, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
	}
	c.SetPingHandler(func(appData string) error {
		return c.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(20*time.Second))
	})
	time_format := "2006-01-02 15:04:05"
	return &Client{host, name, time_format, c}
}

func (client *Client) Close() {
	client.ws.Close()
}

func (client *Client) Read() string {
	_, message, err := client.ws.ReadMessage()
	if err != nil {
		log.Println(err)
	} else {
		return string(message)
	}
	return ""
}

func (client *Client) Write(text string) {
	if err := client.ws.WriteMessage(1, []byte("["+time.Now().Format(client.time_format)+"] "+client.name+": "+text)); err != nil {
		log.Println(err)
	}
}
