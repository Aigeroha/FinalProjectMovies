package repository

import (
	"context"
	"database/sql"
	"final-project/internal/database"
	"final-project/internal/models"
	"strconv"
	"time"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{db: database.DB}
}


func (r *TicketRepository) GetWalletBalance(ctx context.Context, customerID int) (float64, error) {
	var balance float64
	query := "SELECT balance FROM wallets WHERE customer_id = $1"
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&balance)
	return balance, err
}


func (r *TicketRepository) IsSeatTaken(ctx context.Context, scheduleID, seatID int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM tickets WHERE schedule_id = $1 AND seat_id = $2 AND status = 'Куплено'"
	err := r.db.QueryRowContext(ctx, query, scheduleID, seatID).Scan(&count)
	return count > 0, err
}


func (r *TicketRepository) BuyTicketTx(ctx context.Context, t *models.Ticket, price float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	
	walletQuery := "UPDATE wallets SET balance = balance - $1 WHERE customer_id = $2"
	_, err = tx.ExecContext(ctx, walletQuery, price, t.CustomerID)
	if err != nil {
		return err
	}

	
	ticketQuery := `
		INSERT INTO tickets (schedule_id, seat_id, customer_id, ticket_type, status) 
		VALUES ($1, $2, $3, $4, 'Куплено') RETURNING ticket_id`
	err = tx.QueryRowContext(ctx, ticketQuery, t.ScheduleID, t.SeatID, t.CustomerID, t.TicketType).Scan(&t.ID)
	if err != nil {
		return err
	}

	t.Status = "Куплено"
	return tx.Commit()
}


func (r *TicketRepository) CancelAndRefundTicketTx(ctx context.Context, ticketID, customerID int, price float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	
	_, err = tx.ExecContext(ctx, "UPDATE tickets SET status = 'Отмена' WHERE ticket_id = $1", ticketID)
	if err != nil {
		return err
	}


	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1 WHERE customer_id = $2", price, customerID)
	if err != nil {
		return err
	}

	
	refundQuery := "INSERT INTO refunds (ticket_id, refund_date, amount) VALUES ($1, $2, $3)"
	_, err = tx.ExecContext(ctx, refundQuery, ticketID, time.Now(), price)
	if err != nil {
		return err
	}

	return tx.Commit()
}


func (r *TicketRepository) GetByID(ctx context.Context, ticketID int) (*models.Ticket, error) {
	var t models.Ticket
	query := "SELECT ticket_id, schedule_id, seat_id, customer_id, ticket_type, status FROM tickets WHERE ticket_id = $1"
	err := r.db.QueryRowContext(ctx, query, ticketID).Scan(&t.ID, &t.ScheduleID, &t.SeatID, &t.CustomerID, &t.TicketType, &t.Status)
	return &t, err
}


func (r *TicketRepository) GetFilteredTickets(ctx context.Context, filters map[string]string) ([]models.TicketView, error) {
	query := `
		SELECT 
			t.ticket_id, m.title, s.session_date, s.session_time, s.hall_id, 
			se.seat_number, c.name, t.ticket_type, t.status,
			CASE 
				WHEN t.ticket_type = 'Взрослый' THEN s.adult_price
				WHEN t.ticket_type = 'Студенческий' THEN s.student_price
				WHEN t.ticket_type = 'Детский' THEN s.child_price
				ELSE 0.0
			END as price
		FROM tickets t
		JOIN schedules s ON t.schedule_id = s.schedule_id
		JOIN movies m ON s.movie_id = m.movie_id
		JOIN seats se ON t.seat_id = se.seat_id
		JOIN customers c ON t.customer_id = c.customer_id
		WHERE 1=1`

	var args []interface{}
	idx := 1

	if v := filters["movie_id"]; v != "" {
		query += " AND m.movie_id = $" + strconv.Itoa(idx)
		args = append(args, v)
		idx++
	}
	if v := filters["movie_title"]; v != "" {
		query += " AND m.title ILIKE $" + strconv.Itoa(idx)
		args = append(args, "%"+v+"%")
		idx++
	}
	if v := filters["ticket_type"]; v != "" {
		query += " AND t.ticket_type = $" + strconv.Itoa(idx)
		args = append(args, v)
		idx++
	}
	if v := filters["time"]; v != "" {
		if len(v) == 5 {
			v = v + ":00"
		}
		query += " AND s.session_time = $" + strconv.Itoa(idx)
		args = append(args, v)
		idx++
	}
	if v := filters["hall_id"]; v != "" {
		query += " AND s.hall_id = $" + strconv.Itoa(idx)
		args = append(args, v)
		idx++
	}
	if v := filters["schedule_id"]; v != "" {
		query += " AND t.schedule_id = $" + strconv.Itoa(idx)
		args = append(args, v)
		idx++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.TicketView
	for rows.Next() {
		var tv models.TicketView
		err := rows.Scan(&tv.TicketID, &tv.MovieTitle, &tv.SessionDate, &tv.SessionTime, &tv.HallID, &tv.SeatNumber, &tv.CustomerName, &tv.TicketType, &tv.Status, &tv.Price)
		if err != nil {
			return nil, err
		}
		list = append(list, tv)
	}
	return list, nil
}
