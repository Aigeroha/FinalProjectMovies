package services

import (
	"context"
	"database/sql"
	"time"

	"final-project/internal/database"
	"final-project/internal/errs"
	"final-project/internal/models"
)

type MovieService struct{}

func NewMovieService() *MovieService {
	return &MovieService{}
}

// 1. GetMovies — Получить вообще все фильмы из базы
func (s *MovieService) GetMovies(ctx context.Context) ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := database.DB.QueryContext(ctx, "SELECT movie_id, title, duration, genre, rating FROM movies")
	if err != nil {
		return nil, errs.ErrInternal
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Duration, &m.Genre, &m.Rating); err != nil {
			return nil, errs.ErrInternal
		}
		movies = append(movies, m)
	}
	return movies, nil
}

// 2. GetMovieByID — Найти один конкретный фильм по его ID
func (s *MovieService) GetMovieByID(ctx context.Context, id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var m models.Movie
	err := database.DB.QueryRowContext(ctx, "SELECT movie_id, title, duration, genre, rating FROM movies WHERE movie_id = $1", id).
		Scan(&m.ID, &m.Title, &m.Duration, &m.Genre, &m.Rating)

	if err == sql.ErrNoRows {
		return nil, errs.ErrMovieNotFound
	}
	if err != nil {
		return nil, errs.ErrInternal
	}
	return &m, nil
}

// 3. CreateMovie — Добавить новый фильм в базу данных
func (s *MovieService) CreateMovie(ctx context.Context, m *models.Movie) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "INSERT INTO movies (title, duration, genre, rating) VALUES ($1, $2, $3, $4) RETURNING movie_id"
	err := database.DB.QueryRowContext(ctx, query, m.Title, m.Duration, m.Genre, m.Rating).Scan(&m.ID)
	if err != nil {
		return errs.ErrInternal
	}
	return nil
}

// 4. DeleteMovie — Удалить фильм из кинотеатра по ID
func (s *MovieService) DeleteMovie(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := database.DB.ExecContext(ctx, "DELETE FROM movies WHERE movie_id = $1", id)
	if err != nil {
		return errs.ErrInternal
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrMovieNotFound
	}
	return nil
}

// 5. UpdateMovie (PUT) — Полное обновление фильма (перезаписываем все поля)
func (s *MovieService) UpdateMovie(ctx context.Context, id int, m *models.Movie) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "UPDATE movies SET title = $1, duration = $2, genre = $3, rating = $4 WHERE movie_id = $5"
	result, err := database.DB.ExecContext(ctx, query, m.Title, m.Duration, m.Genre, m.Rating, id)
	if err != nil {
		return errs.ErrInternal
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrMovieNotFound
	}
	m.ID = id
	return nil
}

// 6. PatchMovie (PATCH) — Частичное обновление фильма (только то, что прислали)
func (s *MovieService) PatchMovie(ctx context.Context, id int, existing *models.Movie, input map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Проверяем, какие поля лежат в map и обновляем только их
	if val, ok := input["title"].(string); ok {
		existing.Title = val
	}
	if val, ok := input["duration"].(float64); ok { // Числа из JSON в Go прилетают как float64
		existing.Duration = int(val)
	}
	if val, ok := input["genre"].(string); ok {
		existing.Genre = val
	}
	if val, ok := input["rating"].(float64); ok {
		existing.Rating = val
	}

	query := "UPDATE movies SET title = $1, duration = $2, genre = $3, rating = $4 WHERE movie_id = $5"
	_, err := database.DB.ExecContext(ctx, query, existing.Title, existing.Duration, existing.Genre, existing.Rating, id)
	if err != nil {
		return errs.ErrInternal
	}
	return nil
}

// 7. GetMoviePaginated — Пагинация (по 5 фильмов на одну страницу)
func (s *MovieService) GetMoviePaginated(ctx context.Context, page int) ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	limit := 5
	offset := (page - 1) * limit

	query := "SELECT movie_id, title, duration, genre, rating FROM movies LIMIT $1 OFFSET $2"
	rows, err := database.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, errs.ErrInternal
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Duration, &m.Genre, &m.Rating); err != nil {
			return nil, errs.ErrInternal
		}
		movies = append(movies, m)
	}
	return movies, nil
}

// 8. GetMoviesFilter — Поиск фильмов по названию, жанру и минимальному рейтингу
func (s *MovieService) GetMoviesFilter(ctx context.Context, title, genre string, minRating float64) ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	searchTitle := "%" + title + "%"
	searchGenre := "%" + genre + "%"

	// Используем ILIKE для поиска без учета регистра букв
	query := "SELECT movie_id, title, duration, genre, rating FROM movies WHERE title ILIKE $1 AND genre ILIKE $2 AND rating >= $3"
	rows, err := database.DB.QueryContext(ctx, query, searchTitle, searchGenre, minRating)
	if err != nil {
		return nil, errs.ErrInternal
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Duration, &m.Genre, &m.Rating); err != nil {
			return nil, errs.ErrInternal
		}
		movies = append(movies, m)
	}
	return movies, nil
}

// 9. GetMovieStats — Аналитика (Сложная статистика по фильмам через Горутины и Каналы)
func (s *MovieService) GetMovieStats(ctx context.Context) (*models.MovieStats, error) {
	// Для расчета аналитики даем контексту чуть больше времени — 7 секунд
	ctx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()

	// Используем канал для безопасной передачи данных из горутины
	statsChan := make(chan *models.MovieStats, 1)

	// Запускаем расчет параллельно в фоновом потоке
	go func() {
		var stats models.MovieStats
		query := `
			SELECT 
				COUNT(*), 
				COALESCE(AVG(duration), 0), MAX(duration), MIN(duration),
				COALESCE(AVG(rating), 0), MAX(rating), MIN(rating) 
			FROM movies`

		err := database.DB.QueryRowContext(ctx, query).Scan(
			&stats.TotalMovies,
			&stats.Duration.Average,
			&stats.Duration.Max,
			&stats.Duration.Min,
			&stats.Rating.Average,
			&stats.Rating.Max,
			&stats.Rating.Min,
		)
		if err != nil {
			statsChan <- nil
			return
		}
		statsChan <- &stats
	}()

	// Конкурентный выбор с помощью оператора select
	select {
	case <-ctx.Done():
		// Если пользователь закрыл страницу или время вышло
		return nil, errs.ErrTimeout
	case stats := <-statsChan:
		if stats == nil {
			return nil, errs.ErrInternal
		}
		return stats, nil
	}
}
