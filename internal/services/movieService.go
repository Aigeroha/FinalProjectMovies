package services

import (
	"context"
	"final-project/internal/errs"
	"final-project/internal/models"
	"final-project/internal/repository"
)

type MovieService struct {
	repo *repository.MovieRepository
}

func NewMovieService(r *repository.MovieRepository) *MovieService {
	return &MovieService{repo: r}
}

func (s *MovieService) GetMovies(ctx context.Context) ([]models.Movie, error) {
	movies, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, errs.ErrInternal
	}
	return movies, nil
}

func (s *MovieService) GetMovieByID(ctx context.Context, id int) (*models.Movie, error) {
	movie, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) CreateMovie(ctx context.Context, m *models.Movie) error {
	if m.Rating < 0 || m.Rating > 10 {
		return errs.New("рейтинг должен быть от 0 до 10", 400)
	}

	err := s.repo.Create(ctx, m)
	if err != nil {
		return errs.ErrInternal
	}
	return nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, id int) error {
	rowsAffected, err := s.repo.Delete(ctx, id)
	if err != nil {
		return errs.ErrInternal
	}
	if rowsAffected == 0 {
		return errs.ErrMovieNotFound
	}
	return nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, id int, m *models.Movie) error {
	rowsAffected, err := s.repo.Update(ctx, id, m)
	if err != nil {
		return errs.ErrInternal
	}
	if rowsAffected == 0 {
		return errs.ErrMovieNotFound
	}
	m.ID = id
	return nil
}

func (s *MovieService) PatchMovie(ctx context.Context, id int, existing *models.Movie, input map[string]interface{}) error {
	if val, ok := input["title"].(string); ok {
		existing.Title = val
	}
	if val, ok := input["duration"].(float64); ok {
		existing.Duration = int(val)
	}
	if val, ok := input["genre"].(string); ok {
		existing.Genre = val
	}
	if val, ok := input["rating"].(float64); ok {
		existing.Rating = val
	}

	rowsAffected, err := s.repo.Update(ctx, id, existing)
	if err != nil {
		return errs.ErrInternal
	}
	if rowsAffected == 0 {
		return errs.ErrMovieNotFound
	}
	return nil
}

func (s *MovieService) GetMoviePaginated(ctx context.Context, page int) ([]models.Movie, error) {
	limit := 5
	offset := (page - 1) * limit

	movies, err := s.repo.GetPaginated(ctx, limit, offset)
	if err != nil {
		return nil, errs.ErrInternal
	}
	return movies, nil
}

func (s *MovieService) GetMoviesFilter(ctx context.Context, title, genre string, minRating float64) ([]models.Movie, error) {
	searchTitle := "%" + title + "%"
	searchGenre := "%" + genre + "%"

	movies, err := s.repo.GetFiltered(ctx, searchTitle, searchGenre, minRating)
	if err != nil {
		return nil, errs.ErrInternal
	}
	return movies, nil
}

func (s *MovieService) GetMovieStats(ctx context.Context) (*models.MovieStats, error) {
	statsChan := make(chan *models.MovieStats, 1)

	go func() {
		stats, err := s.repo.GetStatsRaw(ctx)
		if err != nil {
			statsChan <- nil
			return
		}
		statsChan <- stats
	}()

	select {
	case <-ctx.Done():
		return nil, errs.ErrTimeout
	case stats := <-statsChan:
		if stats == nil {
			return nil, errs.ErrInternal
		}
		return stats, nil
	}
}
