package repository

import (
	"context"
	"database/sql"
	"errors"
	"final-project/internal/database"
	"final-project/internal/errs"
	"final-project/internal/models"
	"strconv"
)

type ScheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository() *ScheduleRepository {
	return &ScheduleRepository{db: database.DB}
}


func (r *ScheduleRepository) GetAllDetailed(ctx context.Context) ([]models.ScheduleView, error) {
	query := `SELECT schedule_id, movie_title, session_date, session_time, hall_id, adult_price, student_price, child_price 
	          FROM view_readable_schedules`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ScheduleView
	for rows.Next() {
		var sv models.ScheduleView
		err := rows.Scan(
			&sv.ID,
			&sv.MovieTitle,
			&sv.SessionDate,
			&sv.SessionTime,
			&sv.HallID,
			&sv.AdultPrice,
			&sv.StudentPrice,
			&sv.ChildPrice,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, sv)
	}
	return list, nil
}


func (r *ScheduleRepository) GetPaginated(ctx context.Context, limit, offset int) ([]models.ScheduleView, error) {
	query := `SELECT schedule_id, movie_title, session_date, session_time, hall_id, adult_price, student_price, child_price 
	          FROM view_readable_schedules 
	          LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ScheduleView
	for rows.Next() {
		var sv models.ScheduleView
		err := rows.Scan(
			&sv.ID,
			&sv.MovieTitle,
			&sv.SessionDate,
			&sv.SessionTime,
			&sv.HallID,
			&sv.AdultPrice,
			&sv.StudentPrice,
			&sv.ChildPrice,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, sv)
	}
	return list, nil
}


func (r *ScheduleRepository) GetFiltered(ctx context.Context, timeSlot, hallStr, movieTitle string) ([]models.ScheduleView, error) {
	query := `SELECT schedule_id, movie_title, session_date, session_time, hall_id, adult_price, student_price, child_price 
	          FROM view_readable_schedules WHERE 1=1`
	var args []interface{}
	placeholderIdx := 1

	if timeSlot != "" {
		query += " AND session_time = $" + strconv.Itoa(placeholderIdx)
		args = append(args, timeSlot)
		placeholderIdx++
	}
	if hallStr != "" {
		
		hallID, err := strconv.Atoi(hallStr)
		if err == nil {
			query += " AND hall_id = $" + strconv.Itoa(placeholderIdx)
			args = append(args, hallID)
			placeholderIdx++
		}
	}
	if movieTitle != "" {
		query += " AND movie_title ILIKE $" + strconv.Itoa(placeholderIdx)
		args = append(args, "%"+movieTitle+"%")
		placeholderIdx++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ScheduleView
	for rows.Next() {
		var sv models.ScheduleView
		err := rows.Scan(
			&sv.ID,
			&sv.MovieTitle,
			&sv.SessionDate,
			&sv.SessionTime,
			&sv.HallID,
			&sv.AdultPrice,
			&sv.StudentPrice,
			&sv.ChildPrice,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, sv)
	}
	return list, nil
}


func (r *ScheduleRepository) CheckSlotBusy(ctx context.Context, date, timeSlot string, hallID, currentID int) (bool, error) {
	query := "SELECT COUNT(*) FROM schedules WHERE session_date = $1 AND session_time = $2 AND hall_id = $3"
	var args []interface{}
	args = append(args, date, timeSlot, hallID)

	
	if currentID > 0 {
		query += " AND schedule_id != $4"
		args = append(args, currentID)
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}


func (r *ScheduleRepository) Create(ctx context.Context, s *models.Schedule) error {
	query := `INSERT INTO schedules (movie_id, hall_id, session_date, session_time, adult_price, student_price, child_price) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING schedule_id, created_at`
	return r.db.QueryRowContext(ctx, query, s.MovieID, s.HallID, s.SessionDate, s.SessionTime, s.AdultPrice, s.StudentPrice, s.ChildPrice).
		Scan(&s.ID, &s.CreatedAt)
}


func (r *ScheduleRepository) Update(ctx context.Context, id int, s *models.Schedule) (int64, error) {
	query := `UPDATE schedules 
	          SET movie_id = $1, hall_id = $2, session_date = $3, session_time = $4, adult_price = $5, student_price = $6, child_price = $7 
	          WHERE schedule_id = $8`
	result, err := r.db.ExecContext(ctx, query, s.MovieID, s.HallID, s.SessionDate, s.SessionTime, s.AdultPrice, s.StudentPrice, s.ChildPrice, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}


func (r *ScheduleRepository) Delete(ctx context.Context, id int) (int64, error) {
	result, err := r.db.ExecContext(ctx, "DELETE FROM schedules WHERE schedule_id = $1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}


func (r *ScheduleRepository) GetByID(ctx context.Context, id int) (*models.Schedule, error) {
	var s models.Schedule
	query := `SELECT schedule_id, movie_id, hall_id, session_date, session_time, adult_price, student_price, child_price, created_at 
	          FROM schedules WHERE schedule_id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.MovieID,
		&s.HallID,
		&s.SessionDate,
		&s.SessionTime,
		&s.AdultPrice,
		&s.StudentPrice,
		&s.ChildPrice,
		&s.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrScheduleNotFound
	}
	return &s, err
}
