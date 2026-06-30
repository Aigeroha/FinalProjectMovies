package models

import (
	"encoding/json"
	"time"
)

// Schedule используется для CREATE, UPDATE, PATCH (работает через ID)
type Schedule struct {
	ID        int       `json:"id"`
	MovieID   int       `json:"movie_id"`
	HallID    int       `json:"hall_id"`
	Date      string    `json:"date"` // В Postman: "ДД-ММ-ГГГГ" -> В базе: "ГГГГ-ММ-ДД"
	Time      string    `json:"time"` // В Postman: "ЧЧ:ММ" -> В базе: "ЧЧ:ММ:00"
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type ScheduleView struct {
	ID         int     `json:"schedule_id"`
	MovieTitle string  `json:"movie_title"`
	HallName   string  `json:"hall_name"`
	Date       string  `json:"date"` // Будет выводиться как "ДД-ММ-ГГГГ"
	Time       string  `json:"time"` // Будет выводиться как "ЧЧ:ММ"
	Price      float64 `json:"price"`
}


func (s *Schedule) UnmarshalJSON(data []byte) error {
	type Alias Schedule
	aux := &struct{ *Alias }{Alias: (*Alias)(s)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if s.Date != "" {
		t, err := time.Parse("02-01-2006", s.Date)
		if err == nil {
			s.Date = t.Format("2006-01-02")
		}
	}
	if s.Time != "" && len(s.Time) == 5 {
		s.Time = s.Time + ":00"
	}
	return nil
}


func (s *Schedule) MarshalJSON() ([]byte, error) {
	type Alias Schedule
	return json.Marshal(&struct {
		*Alias
		Date string `json:"date"`
		Time string `json:"time"`
	}{
		Alias: (*Alias)(s),
		Date:  formatToDisplayDate(s.Date),
		Time:  formatToDisplayTime(s.Time),
	})
}

func (sv *ScheduleView) MarshalJSON() ([]byte, error) {
	type Alias ScheduleView
	return json.Marshal(&struct {
		*Alias
		Date string `json:"date"`
		Time string `json:"time"`
	}{
		Alias: (*Alias)(sv),
		Date:  formatToDisplayDate(sv.Date),
		Time:  formatToDisplayTime(sv.Time),
	})
}


func formatToDisplayDate(dbDate string) string {
	t, err := time.Parse("2006-01-02", dbDate)
	if err != nil {
		return dbDate
	}
	return t.Format("02-01-2006")
}

func formatToDisplayTime(dbTime string) string {
	if len(dbTime) >= 5 {
		return dbTime[:5]
	}
	return dbTime
}
