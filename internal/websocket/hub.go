package websocket

import (
	"log"
	"sync"

	fiberws "github.com/gofiber/contrib/websocket"
)

type Hub struct {
	clients sync.Map
}

var AppHub = &Hub{}

func (h *Hub) Register(conn *fiberws.Conn) {
	h.clients.Store(conn, true)
	log.Printf("Новый клиент подключился: %s", conn.RemoteAddr())
}

func (h *Hub) Unregister(conn *fiberws.Conn) {
	h.clients.Delete(conn)
	log.Printf("Клиент отключился: %s", conn.RemoteAddr())
}

func (h *Hub) Broadcast(message []byte) {
	h.clients.Range(func(key, value interface{}) bool {
		conn := key.(*fiberws.Conn)
		err := conn.WriteMessage(fiberws.TextMessage, message)
		if err != nil {
			log.Printf("Ошибка отправки: %v", err)
			conn.Close()
			h.Unregister(conn)
		}
		return true
	})
}