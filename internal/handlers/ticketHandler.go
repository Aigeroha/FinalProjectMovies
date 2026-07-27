package handlers

import (
	"context"
	"strconv"
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

	val := c.Locals("customer_id")
	if val == nil {
		return responses.Error(c, 401, "unauthorized")
	}
	var customerID int
	switch id := val.(type) {
	case int:
		customerID = id
	case float64:
		customerID = int(id)
	case uint:
		customerID = int(id)
	default:
		return responses.Error(c, 401, "неверный формат ID пользователя в токенe")
	}

	t.CustomerID = customerID

	if err := h.service.CreateTicket(ctx, &t); err != nil {
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, 201, t)
}

func (h *TicketHandler) RefundTicket(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var input struct {
		TicketID int `json:"ticket_id"`
	}
	if err := c.Bind().Body(&input); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	customerID, ok := c.Locals("customer_id").(int)
	if !ok {
		return responses.Error(c, 401, "unauthorized")
	}

	if err := h.service.RefundTicket(ctx, input.TicketID, customerID); err != nil {
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

func (h *TicketHandler) GetTicketByID(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return responses.Error(c, 400, "неверный id")
	}

	ticket, err := h.service.GetTicketByID(ctx, id)
	if err != nil {
		return responses.Error(c, 404, err.Error())
	}

	return responses.Success(c, 200, ticket)
}
