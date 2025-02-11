package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte("\n")
	space   = []byte(" ")
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (client *Client) readPump() {
	defer func() {
		client.Hub.Unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	// client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// client.Conn.SetPongHandler(func(appData string) error {
	// 	fmt.Printf("ping! =%v", appData)
	// 	return client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// })

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway)
		}
		client.Hub.Broadcast <- message
	}
}

func (client *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			wc, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			wc.Write(message)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				wc.Write(newline)
				wc.Write(<-client.Send)
			}

			if err := wc.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.TextMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.readPump()
	go client.writePump()

}
