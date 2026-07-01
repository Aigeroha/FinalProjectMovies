package handlers

import (
	"context"
	"strconv"
	"time"

	"final-project/internal/models"
	"final-project/internal/responses"
	"final-project/internal/services"
	"final-project/internal/auth"
	"github.com/gofiber/fiber/v3"
)

type CustomerHandler struct {
	service *services.CustomerService
}

func NewCustomerHandler(s *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: s}
}


func (h *CustomerHandler) Register(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var input models.RegisterInput
	if err := c.Bind().Body(&input); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	if err := h.service.Register(ctx, input); err != nil {
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, 201, map[string]string{"message": "регистрация успешна"})
}


func (h *CustomerHandler) Login(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var input models.LoginInput
	if err := c.Bind().Body(&input); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	customer, err := h.service.Login(ctx, input)
	if err != nil {
		return responses.Error(c, 401, err.Error())
	}

	
	token, err := auth.GenerateToken(customer.ID)
	if err != nil {
		return responses.Error(c, 500, "could not generate token")
	}

	
	return responses.Success(c, 200, map[string]string{
		"token": token,
	})
}


func (h *CustomerHandler) GetProfile(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid customer id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profile, err := h.service.GetCustomerProfile(ctx, id)
	if err != nil {
		return responses.Error(c, 404, err.Error())
	}

	return responses.Success(c, 200, profile)
}


func (h *CustomerHandler) PatchCustomer(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid customer id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var input map[string]interface{}
	if err := c.Bind().Body(&input); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	if err := h.service.PatchCustomer(ctx, id, input); err != nil {
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, 200, map[string]string{"message": "данные обновлены"})
}


func (h *CustomerHandler) TopUpWallet(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid customer id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var input struct {
		Amount int `json:"amount"`
	}
	if err := c.Bind().Body(&input); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	if err := h.service.TopUpWallet(ctx, id, input.Amount); err != nil {
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, 200, map[string]string{"message": "кошелек успешно пополнен"})
}


func (h *CustomerHandler) DeleteCustomer(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid customer id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.service.DeleteCustomer(ctx, id); err != nil {
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, 200, map[string]string{"message": "аккаунт и кошелек успешно удалены"})
}


func (h *CustomerHandler) GetAllCustomers(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return responses.Error(c, 400, "invalid page parameter")
	}

	list, err := h.service.GetAllCustomersPaginated(ctx, page)
	if err != nil {
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, list)
}
