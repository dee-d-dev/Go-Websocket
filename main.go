package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	clients := make([]*websocket.Conn, 0)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		clients = append(clients, conn)

		for {
			remoteAddr := conn.RemoteAddr().String()
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			for _, client := range clients {
				payload, err := json.Marshal(map[string]string{
					"remoteAddress": remoteAddr,
					"message":       string(msg),
				})

				if err != nil {
					fmt.Println("Error marshalling json")
				}
				err = client.WriteMessage(msgType, payload)

				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	fmt.Println("starting server on port 5000")

	http.ListenAndServe(":5000", nil)
}
