package models

type Ticket struct {
	ID         int    `json:"ticket_id"`
	ScheduleID int    `json:"schedule_id"`
	SeatID     int    `json:"seat_id"`
	CustomerID int    `json:"customer_id"`
	TicketType string `json:"ticket_type"`
	Status     string `json:"status"`
}

type TicketView struct {
	TicketID     int     `json:"ticket_id"`
	MovieTitle   string  `json:"movie_title"`
	SessionDate  string  `json:"session_date"`
	SessionTime  string  `json:"session_time"`
	HallID       int     `json:"hall_id"`
	SeatNumber   int     `json:"seat_number"`
	CustomerName string  `json:"customer_name"`
	TicketType   string  `json:"ticket_type"`
	Price        float64 `json:"price"`
	Status       string  `json:"status"`
}
