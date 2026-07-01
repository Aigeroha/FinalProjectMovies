package handlers

import (
	"log"

	"final-project/internal/websocket"
	fiberws "github.com/gofiber/contrib/websocket"
)


type SeatSelectionMessage struct {
	ScheduleID int  `json:"schedule_id"`
	SeatID     int  `json:"seat_id"`
	IsBooked   bool `json:"is_booked"`
}


func HandleWebSocket(c *fiberws.Conn) {
	
	websocket.AppHub.Register(c)
	
	defer func() {
		websocket.AppHub.Unregister(c)
		c.Close()
	}()

	for {
		
		messageType, msgBytes, err := c.ReadMessage()
		if err != nil {
			log.Printf("Ошибка чтения сообщения: %v", err)
			break
		}

		
		if messageType == fiberws.TextMessage {
			log.Printf("Получено сообщение от клиента, отправляем в Broadcast...")
			websocket.AppHub.Broadcast(msgBytes)
		}
	}
}