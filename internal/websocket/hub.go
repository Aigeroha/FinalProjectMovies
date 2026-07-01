package websocket

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)


type Hub struct {
	
	clients sync.Map
}

var AppHub = &Hub{}


func (h *Hub) Register(conn *websocket.Conn) {
	h.clients.Store(conn, true)
	log.Printf("Новый клиент подключился: %s", conn.RemoteAddr())
}


func (h *Hub) Unregister(conn *websocket.Conn) {
	h.clients.Delete(conn)
	log.Printf("Клиент отключился: %s", conn.RemoteAddr())
}


func (h *Hub) Broadcast(message []byte) {
	h.clients.Range(func(key, value interface{}) bool {
		conn := key.(*websocket.Conn)
		
		
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Ошибка отправки сообщения клиенту %s: %v", conn.RemoteAddr(), err)
			conn.Close()
			h.Unregister(conn)
		}
		return true 
	})
}