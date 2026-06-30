package handlers

import (
	"context"
	"time"

	"final-project/internal/models"
	"final-project/internal/responses"
	"final-project/internal/services"

	"github.com/gofiber/fiber/v3"
)

type TicketHandler struct {
	service *services.TicketService
}

func NewTicketHandler(s *services.TicketService) *TicketHandler {
	return &TicketHandler{service: s}
}


func (h *TicketHandler) BuyTicket(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var t models.Ticket
	if err := c.Bind().Body(&t); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	if err := h.service.CreateTicket(ctx, &t); err != nil {
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, 201, t)
}


func (h *TicketHandler) RefundTicket(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	
	var input struct {
		TicketID   int `json:"ticket_id"`
		CustomerID int `json:"customer_id"`
	}
	if err := c.Bind().Body(&input); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	if err := h.service.RefundTicket(ctx, input.TicketID, input.CustomerID); err != nil {
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, 200, map[string]string{"message": "билет успешно отменен, средства возвращены"})
}


func (h *TicketHandler) GetTickets(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := map[string]string{
		"movie_id":    c.Query("movie_id"),
		"movie_title": c.Query("movie_title"),
		"ticket_type": c.Query("ticket_type"), 
		"time":        c.Query("time"),
		"hall_id":     c.Query("hall_id"),
		"schedule_id": c.Query("schedule_id"),
	}

	list, err := h.service.GetTickets(ctx, filters)
	if err != nil {
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, list)
}
