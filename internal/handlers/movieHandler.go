package handlers

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"final-project/internal/errs"
	"final-project/internal/models"
	"final-project/internal/responses"
	"final-project/internal/services"

	"github.com/gofiber/fiber/v3"
)

type MovieHandler struct {
	service *services.MovieService
}

func NewMovieHandler(s *services.MovieService) *MovieHandler {
	return &MovieHandler{service: s}
}

func (h *MovieHandler) GetAllMovies(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос GetMovies обработан за: %v", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movies, err := h.service.GetMovies(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return responses.Error(c, 408, "request timeout")
		}
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, movies)
}

func (h *MovieHandler) GetMoviesFilter(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос GetMoviesFilter обработан за: %v", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	title := c.Query("title")
	genre := c.Query("genre")
	ratingStr := c.Query("rating")

	var minRating float64
	if ratingStr != "" {
		var err error
		minRating, err = strconv.ParseFloat(ratingStr, 64)
		if err != nil {
			return responses.Error(c, 400, "invalid rating parameter")
		}
	}

	movies, err := h.service.GetMoviesFilter(ctx, title, genre, minRating)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return responses.Error(c, 408, "request timeout")
		}
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, movies)
}

func (h *MovieHandler) GetMoviesPaginated(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос GetMoviesPaginated processed in: %v", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return responses.Error(c, 400, "invalid page parameter")
	}

	movies, err := h.service.GetMoviePaginated(ctx, page)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return responses.Error(c, 408, "request timeout")
		}
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, movies)
}

func (h *MovieHandler) GetMovieByID(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос GetMovieByID обработан за: %v", time.Since(start))
	}()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid movie id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movie, err := h.service.GetMovieByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return responses.Error(c, 408, "request timeout")
		}
		if errors.Is(err, errs.ErrMovieNotFound) {
			return responses.Error(c, 404, "movie not found")
		}
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, movie)
}

func (h *MovieHandler) CreateMovie(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос CreateMovie обработан за: %v", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var movie models.Movie
	if err := c.Bind().Body(&movie); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	err := h.service.CreateMovie(ctx, &movie)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return responses.Error(c, 408, "request timeout")
		}
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 201, movie)
}

func (h *MovieHandler) UpdateMovie(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос UpdateMovie обработан за: %v", time.Since(start))
	}()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid movie id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var movie models.Movie
	if err := c.Bind().Body(&movie); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	err = h.service.UpdateMovie(ctx, id, &movie)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return responses.Error(c, 408, "request timeout")
		}
		if errors.Is(err, errs.ErrMovieNotFound) {
			return responses.Error(c, 404, "movie not found")
		}
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, movie)
}

func (h *MovieHandler) PatchMovie(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос PatchMovie обработан за: %v", time.Since(start))
	}()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid movie id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var input map[string]interface{}
	if err := c.Bind().Body(&input); err != nil {
		return responses.Error(c, 400, "bad request")
	}

	existing, err := h.service.GetMovieByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return responses.Error(c, 408, "request timeout")
		}
		if errors.Is(err, errs.ErrMovieNotFound) {
			return responses.Error(c, 404, "movie not found")
		}
		return responses.Error(c, 500, "internal server error")
	}

	err = h.service.PatchMovie(ctx, id, existing, input)
	if err != nil {
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, existing)
}

func (h *MovieHandler) DeleteMovie(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос DeleteMovie обработан за: %v", time.Since(start))
	}()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return responses.Error(c, 400, "invalid movie id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.service.DeleteMovie(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return responses.Error(c, 408, "request timeout")
		}
		if errors.Is(err, errs.ErrMovieNotFound) {
			return responses.Error(c, 404, "movie not found")
		}
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, map[string]string{"message": "movie deleted"})
}

func (h *MovieHandler) GetMovieStats(c fiber.Ctx) error {
	start := time.Now()
	defer func() {
		log.Printf("Запрос GetMovieStats обработан за: %v", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stats, err := h.service.GetMovieStats(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, errs.ErrTimeout) {
			return responses.Error(c, 408, "request timeout")
		}
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, 200, stats)
}
