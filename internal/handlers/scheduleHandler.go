package handlers

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"final-project/internal/errs"
	"final-project/internal/models"
	"final-project/internal/responses"
	"final-project/internal/services"
)

type ScheduleHandler struct {
	service *services.ScheduleService
}

func NewScheduleHandler(s *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: s}
}

func (h *ScheduleHandler) GetSchedules(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос GetSchedules обработан за: %v", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	timeSlot := c.Query("time")
	hallName := c.Query("hall")
	movieTitle := c.Query("movie")

	
	if timeSlot != "" || hallName != "" || movieTitle != "" {
		list, err := h.service.GetSchedulesFilter(ctx, timeSlot, hallName, movieTitle)
		if err != nil {
			return responses.Error(c, 500, "internal server error")
		}
		return responses.Success(c, 200, list)
	}

	
	list, err := h.service.GetSchedules(ctx)
	if err != nil {
		return responses.Error(c, 500, "internal server error")
	}
	return responses.Success(c, 200, list)
}


func (h *ScheduleHandler) GetSchedulesPaginated(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос GetSchedulesPaginated обработан за: %v", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return responses.Error(c, 400, "invalid page parameter")
	}

	list, err := h.service.GetSchedulesPaginated(ctx, page)
	if err != nil {
		return responses.Error(c, 500, "internal server error")
	}
	return responses.Success(c, 200, list)
}


func (h *ScheduleHandler) CreateSchedule(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос CreateSchedule обработан за: %v", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sch models.Schedule
	if err := c.Bind().Body(&sch); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	err := h.service.CreateSchedule(ctx, &sch)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return responses.Error(c, 408, "request timeout")
		}
		
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, 201, sch)
}


func (h *ScheduleHandler) UpdateSchedule(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос UpdateSchedule обработан за: %v", time.Since(start))
	}()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid schedule id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sch models.Schedule
	if err := c.Bind().Body(&sch); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	err = h.service.UpdateSchedule(ctx, id, &sch)
	if err != nil {
		if errors.Is(err, errs.ErrScheduleNotFound) {
			return responses.Error(c, 404, "schedule not found")
		}
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, 200, sch)
}

func (h *ScheduleHandler) PatchSchedule(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос PatchSchedule обработан за: %v", time.Since(start))
	}()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid schedule id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var input map[string]interface{}
	if err := c.Bind().Body(&input); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	existing, err := h.service.GetScheduleByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrScheduleNotFound) {
			return responses.Error(c, 404, "schedule not found")
		}
		return responses.Error(c, 500, "internal server error")
	}

	err = h.service.PatchSchedule(ctx, id, existing, input)
	if err != nil {
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, 200, existing)
}


func (h *ScheduleHandler) DeleteSchedule(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос DeleteSchedule обработан за: %v", time.Since(start))
	}()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid schedule id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.service.DeleteSchedule(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrScheduleNotFound) {
			return responses.Error(c, 404, "schedule not found")
		}
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, map[string]string{"message": "schedule deleted"})
}
