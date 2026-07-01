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

func HandlerWebSocket(c *fiberws.Conn) error { 
	websocket.AppHub.Register(c) // А здесь используем оригинальный websocket (твой пакет)
	
	defer func() {
		websocket.AppHub.Unregister(c)
		c.Close()
	}()

	for {
		// И здесь тоже используем алиас fiberws
		messageType, msgBytes, err := c.ReadMessage()
		if err != nil {
			log.Printf("Ошибка: %v", err)
			break
		}
		// И здесь
		if messageType == fiberws.TextMessage {
			websocket.AppHub.Broadcast(msgBytes)
		}
	}
	return nil
}
