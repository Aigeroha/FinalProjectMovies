package models

import (
	"encoding/json"
	"time"
)

type Schedule struct {
	ID           int       `json:"id"`
	MovieID      int       `json:"movie_id"`
	HallID       int       `json:"hall_id"`
	SessionDate  string    `json:"session_date"`
	SessionTime  string    `json:"session_time"`
	AdultPrice   float64   `json:"adult_price"`
	StudentPrice float64   `json:"student_price"`
	ChildPrice   float64   `json:"child_price"`
	CreatedAt    time.Time `json:"created_at"`
}

type ScheduleView struct {
	ID           int     `json:"schedule_id"`
	MovieTitle   string  `json:"movie_title"`
	SessionDate  string  `json:"session_date"`
	SessionTime  string  `json:"session_time"`
	HallID       int     `json:"hall_id"`
	AdultPrice   float64 `json:"adult_price"`
	StudentPrice float64 `json:"student_price"`
	ChildPrice   float64 `json:"child_price"`
}

func (s *Schedule) UnmarshalJSON(data []byte) error {
	type Alias Schedule
	aux := &struct{ *Alias }{Alias: (*Alias)(s)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if s.SessionDate != "" {
		t, err := time.Parse("02-01-2006", s.SessionDate)
		if err == nil {
			s.SessionDate = t.Format("2006-01-02")
		}
	}
	if s.SessionTime != "" && len(s.SessionTime) == 5 {
		s.SessionTime = s.SessionTime + ":00"
	}
	return nil
}

func (s *Schedule) MarshalJSON() ([]byte, error) {
	type Alias Schedule
	return json.Marshal(&struct {
		*Alias
		SessionDate string `json:"session_date"`
		SessionTime string `json:"session_time"`
	}{
		Alias:       (*Alias)(s),
		SessionDate: formatToDisplayDate(s.SessionDate),
		SessionTime: formatToDisplayTime(s.SessionTime),
	})
}

// MarshalJSON для красивого вывода во Вьюшке
func (sv *ScheduleView) MarshalJSON() ([]byte, error) {
	type Alias ScheduleView
	return json.Marshal(&struct {
		*Alias
		SessionDate string `json:"session_date"`
		SessionTime string `json:"session_time"`
	}{
		Alias:       (*Alias)(sv),
		SessionDate: formatToDisplayDate(sv.SessionDate),
		SessionTime: formatToDisplayTime(sv.SessionTime),
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
