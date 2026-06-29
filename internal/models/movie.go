package models

type Movie struct {
	ID       int     `json:"movie_id"`
	Title    string  `json:"title"`
	Duration int     `json:"duration"`
	Genre    string  `json:"genre"`
	Rating   float64 `json:"rating"`
}

type MovieStats struct {
	TotalMovies int           `json:"total_movies"`
	Duration    DurationStats `json:"duration"`
	Rating      RatingStats   `json:"rating"`
}

// DurationStats описывает статистику по длительности фильмов
type DurationStats struct {
	Average float64 `json:"average"`
	Max     int     `json:"max"`
	Min     int     `json:"min"`
}

// RatingStats описывает статистику по рейтингам фильмов
type RatingStats struct {
	Average float64 `json:"average"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}
