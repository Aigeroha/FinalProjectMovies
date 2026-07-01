package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"final-project/internal/database"
	"final-project/internal/errs" 
	"final-project/internal/models"
	"final-project/internal/repository"

	"github.com/redis/go-redis/v9"
)

type MovieService struct {
	repo *repository.MovieRepository
}

func NewMovieService(r *repository.MovieRepository) *MovieService {
	return &MovieService{repo: r}
}


func (s *MovieService) GetAllMovies(ctx context.Context) ([]models.Movie, error) {
	cacheKey := "movies:all"

	
	val, err := database.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedMovies []models.Movie
		if err := json.Unmarshal([]byte(val), &cachedMovies); err == nil {
			log.Println(" [Redis] Данные успешно отданы из кэша!")
			return cachedMovies, nil
		}
	} else if err != redis.Nil {
		log.Printf("Предупреждение: ошибка чтения из Redis: %v", err)
	}

	
	log.Println(" [Postgres] Кэш пуст. Запрос отправлен в базу данных...")
	movies, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	
	moviesJSON, err := json.Marshal(movies)
	if err == nil {
		database.RDB.Set(ctx, cacheKey, moviesJSON, 10*time.Minute)
		log.Println(" [Redis] Свежие данные успешно закэшированы на 10 минут.")
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

	
	database.RDB.Del(ctx, "movies:all")
	log.Println(" [Redis] Кэш очищен из-за создания нового фильма.")
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

	
	database.RDB.Del(ctx, "movies:all")
	log.Println(" [Redis] Кэш очищен из-за удаления фильма.")
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

	
	database.RDB.Del(ctx, "movies:all")
	log.Println(" [Redis] Кэш очищен из-за обновления фильма.")
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

	
	database.RDB.Del(ctx, "movies:all")
	log.Println(" [Redis] Кэш очищен из-за частичного изменения фильма (PATCH).")
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
