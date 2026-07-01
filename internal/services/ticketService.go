package services

import (
	"context"
	"time"

	"final-project/internal/errs"
	"final-project/internal/models"
	"final-project/internal/repository"
)

type TicketService struct {
	repo         *repository.TicketRepository
	scheduleRepo *repository.ScheduleRepository 
}

func NewTicketService(r *repository.TicketRepository, s *repository.ScheduleRepository) *TicketService {
	return &TicketService{repo: r, scheduleRepo: s}
}


func (s *TicketService) CreateTicket(ctx context.Context, t *models.Ticket) error {
	
	sch, err := s.scheduleRepo.GetByID(ctx, t.ScheduleID)
	if err != nil {
		return errs.New("сеанс не найден", 404)
	}

	
	sessionDateTimeStr := sch.SessionDate + " " + sch.SessionTime
	sessionTime, err := time.Parse("2006-01-02 15:04:05", sessionDateTimeStr)
	if err != nil {
		
		sessionTime, err = time.Parse("2006-01-02 15:04", sessionDateTimeStr)
	}

	if err == nil && time.Now().After(sessionTime) {
		return errs.New("нельзя купить билет на уже начавшийся или прошедший сеанс", 400)
	}

	
	var price float64
	switch t.TicketType {
	case "Взрослый":
		price = sch.AdultPrice
	case "Студенческий":
		price = sch.StudentPrice
	case "Детский":
		price = sch.ChildPrice
	default:
		return errs.New("неверный тип билета. Допустимы: 'Взрослый', 'Студенческий', 'Детский'", 400)
	}

	
	balance, err := s.repo.GetWalletBalance(ctx, t.CustomerID)
	if err != nil {
		return errs.New("кошелек пользователя не найден", 404)
	}
	if balance < price {
		return errs.New("не достаточно средств", 400)
	}

	
	isSeatBusy, err := s.repo.IsSeatTaken(ctx, t.ScheduleID, t.SeatID)
	if err != nil {
		return errs.ErrInternal
	}
	if isSeatBusy {
		return errs.New("это место на выбранный сеанс уже занято", 409)
	}


	err = s.repo.BuyTicketTx(ctx, t, price)
	if err != nil {
		return errs.ErrInternal
	}

	return nil
}


func (s *TicketService) RefundTicket(ctx context.Context, ticketID, customerID int) error {
	
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return errs.New("билет не найден", 404)
	}
	if ticket.CustomerID != customerID {
		return errs.New("этот билет принадлежит другому клиенту", 403)
	}
	if ticket.Status == "Отмена" {
		return errs.New("билет уже был отменен ранее", 400)
	}

	
	sch, err := s.scheduleRepo.GetByID(ctx, ticket.ScheduleID)
	if err != nil {
		return errs.ErrInternal
	}
	sessionDateTimeStr := sch.SessionDate + " " + sch.SessionTime
	sessionTime, _ := time.Parse("2006-01-02 15:04:05", sessionDateTimeStr)

	if time.Until(sessionTime) < 1*time.Hour {
		return errs.New("отмена билета невозможна. До начала сеанса осталось меньше часа", 400)
	}


	var refundAmount float64
	switch ticket.TicketType {
	case "Взрослый":
		refundAmount = sch.AdultPrice
	case "Студенческий":
		refundAmount = sch.StudentPrice
	case "Детский":
		refundAmount = sch.ChildPrice
	}

	
	err = s.repo.CancelAndRefundTicketTx(ctx, ticketID, customerID, refundAmount)
	if err != nil {
		return errs.ErrInternal
	}

	return nil
}


func (s *TicketService) GetTickets(ctx context.Context, filter map[string]string) ([]models.TicketView, error) {
	tickets, err := s.repo.GetFilteredTickets(ctx, filter)
	if err != nil {
		return nil, errs.ErrInternal
	}
	return tickets, nil
}