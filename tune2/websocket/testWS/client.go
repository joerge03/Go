package main

import (
	"bytes"
	"fmt"
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
		fmt.Println("close")
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(appData string) error {
		fmt.Printf("ping! =%v\n", appData)
		return client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		fmt.Println("Read")
		_, m, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseInternalServerErr) {
				log.Panicf("WS Error: %v\n", err)
			}
			return
		}
		message := bytes.TrimSpace(bytes.Replace(m, newline, space, -1))
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
			ws, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			ws.Write(message)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				ws.Write(newline)
				ws.Write(<-client.Send)
			}

			if err := ws.Close(); err != nil {
				return
			}
		case <-ticker.C:
			fmt.Println("ticker")
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err, "upgrader failed")
		return
	}

	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	// fmt.Printf("%+v\n", client.Hub.Clients)

	go client.readPump()
	go client.writePump()

}
