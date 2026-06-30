package repository

import (
	"context"
	"database/sql"
	"final-project/internal/database"
	"final-project/internal/errs"
	"final-project/internal/models"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository() *MovieRepository {
	return &MovieRepository{db: database.DB}
}

func (r *MovieRepository) GetAll(ctx context.Context) ([]models.Movie, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT movie_id, title, duration, genre, rating FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Duration, &m.Genre, &m.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *MovieRepository) GetByID(ctx context.Context, id int) (*models.Movie, error) {
	var m models.Movie
	err := r.db.QueryRowContext(ctx, "SELECT movie_id, title, duration, genre, rating FROM movies WHERE movie_id = $1", id).
		Scan(&m.ID, &m.Title, &m.Duration, &m.Genre, &m.Rating)

	if err == sql.ErrNoRows {
		return nil, errs.ErrMovieNotFound
	}
	return &m, err
}

func (r *MovieRepository) Create(ctx context.Context, m *models.Movie) error {
	query := "INSERT INTO movies (title, duration, genre, rating) VALUES ($1, $2, $3, $4) RETURNING movie_id"
	return r.db.QueryRowContext(ctx, query, m.Title, m.Duration, m.Genre, m.Rating).Scan(&m.ID)
}

func (r *MovieRepository) Update(ctx context.Context, id int, m *models.Movie) (int64, error) {
	query := "UPDATE movies SET title = $1, duration = $2, genre = $3, rating = $4 WHERE movie_id = $5"
	result, err := r.db.ExecContext(ctx, query, m.Title, m.Duration, m.Genre, m.Rating, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *MovieRepository) Delete(ctx context.Context, id int) (int64, error) {
	result, err := r.db.ExecContext(ctx, "DELETE FROM movies WHERE movie_id = $1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *MovieRepository) GetPaginated(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	query := "SELECT movie_id, title, duration, genre, rating FROM movies LIMIT $1 OFFSET $2"
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Duration, &m.Genre, &m.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *MovieRepository) GetFiltered(ctx context.Context, title, genre string, minRating float64) ([]models.Movie, error) {
	query := "SELECT movie_id, title, duration, genre, rating FROM movies WHERE title ILIKE $1 AND genre ILIKE $2 AND rating >= $3"
	rows, err := r.db.QueryContext(ctx, query, title, genre, minRating)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Duration, &m.Genre, &m.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *MovieRepository) GetStatsRaw(ctx context.Context) (*models.MovieStats, error) {
	var stats models.MovieStats
	query := `
		SELECT 
			COUNT(*), 
			COALESCE(AVG(duration), 0), MAX(duration), MIN(duration),
			COALESCE(AVG(rating), 0), MAX(rating), MIN(rating) 
		FROM movies`

	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalMovies,
		&stats.Duration.Average,
		&stats.Duration.Max,
		&stats.Duration.Min,
		&stats.Rating.Average,
		&stats.Rating.Max,
		&stats.Rating.Min,
	)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
