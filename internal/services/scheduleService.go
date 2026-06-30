package services

import (
	"context"
	"final-project/internal/errs"
	"final-project/internal/models"
	"final-project/internal/repository"
)

type ScheduleService struct {
	repo *repository.ScheduleRepository
}

func NewScheduleService(r *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: r}
}


func (s *ScheduleService) GetSchedules(ctx context.Context) ([]models.ScheduleView, error) {
	schedules, err := s.repo.GetAllDetailed(ctx)
	if err != nil {
		return nil, errs.ErrInternal
	}
	return schedules, nil
}


func (s *ScheduleService) GetSchedulesPaginated(ctx context.Context, page int) ([]models.ScheduleView, error) {
	limit := 10
	offset := (page - 1) * limit

	schedules, err := s.repo.GetPaginated(ctx, limit, offset)
	if err != nil {
		return nil, errs.ErrInternal
	}
	return schedules, nil
}


func (s *ScheduleService) GetSchedulesFilter(ctx context.Context, timeSlot, hallName, movieTitle string) ([]models.ScheduleView, error) {
	if timeSlot != "" && len(timeSlot) == 5 {
		timeSlot = timeSlot + ":00"
	}

	schedules, err := s.repo.GetFiltered(ctx, timeSlot, hallName, movieTitle)
	if err != nil {
		return nil, errs.ErrInternal
	}
	return schedules, nil
}


func (s *ScheduleService) CreateSchedule(ctx context.Context, sch *models.Schedule) error {
	
	if !s.isValidTimeSlot(sch.SessionTime) {
		return errs.New("неверный временной слот. Доступные слоты: 12:00, 16:00, 19:30, 22:00, 23:50", 400)
	}

	isBusy, err := s.repo.CheckSlotBusy(ctx, sch.SessionDate, sch.SessionTime, sch.HallID, 0)
	if err != nil {
		return errs.ErrInternal
	}
	if isBusy {
		return errs.New("выбранный зал в это время уже занят другим сеансом", 409)
	}

	err = s.repo.Create(ctx, sch)
	if err != nil {
		return errs.ErrInternal
	}
	return nil
}


func (s *ScheduleService) UpdateSchedule(ctx context.Context, id int, sch *models.Schedule) error {
	
	if !s.isValidTimeSlot(sch.SessionTime) {
		return errs.New("неверный временной слот. Доступные слоты: 12:00, 16:00, 19:30, 22:00, 23:50", 400)
	}

	isBusy, err := s.repo.CheckSlotBusy(ctx, sch.SessionDate, sch.SessionTime, sch.HallID, id)
	if err != nil {
		return errs.ErrInternal
	}
	if isBusy {
		return errs.New("выбранный зал в это время уже занят другим сеансом", 409)
	}

	rowsAffected, err := s.repo.Update(ctx, id, sch)
	if err != nil {
		return errs.ErrInternal
	}
	if rowsAffected == 0 {
		return errs.ErrScheduleNotFound
	}
	sch.ID = id
	return nil
}


func (s *ScheduleService) PatchSchedule(ctx context.Context, id int, existing *models.Schedule, input map[string]interface{}) error {
	if val, ok := input["movie_id"].(float64); ok {
		existing.MovieID = int(val)
	}
	if val, ok := input["hall_id"].(float64); ok {
		existing.HallID = int(val)
	}

	
	if val, ok := input["adult_price"].(float64); ok {
		existing.AdultPrice = val
	}
	if val, ok := input["student_price"].(float64); ok {
		existing.StudentPrice = val
	}
	if val, ok := input["child_price"].(float64); ok {
		existing.ChildPrice = val
	}

	
	if val, ok := input["session_date"].(string); ok {
		existing.SessionDate = val
	}
	if val, ok := input["session_time"].(string); ok {
		existing.SessionTime = val
	}

	if !s.isValidTimeSlot(existing.SessionTime) {
		return errs.New("неверный временной слот. Доступные слоты: 12:00, 16:00, 19:30, 22:00, 23:50", 400)
	}

	isBusy, err := s.repo.CheckSlotBusy(ctx, existing.SessionDate, existing.SessionTime, existing.HallID, id)
	if err != nil {
		return errs.ErrInternal
	}
	if isBusy {
		return errs.New("выбранный зал в это время уже занят другим сеансом", 409)
	}

	rowsAffected, err := s.repo.Update(ctx, id, existing)
	if err != nil {
		return errs.ErrInternal
	}
	if rowsAffected == 0 {
		return errs.ErrScheduleNotFound
	}
	return nil
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, id int) error {
	rowsAffected, err := s.repo.Delete(ctx, id)
	if err != nil {
		return errs.ErrInternal
	}
	if rowsAffected == 0 {
		return errs.ErrScheduleNotFound
	}
	return nil
}

func (s *ScheduleService) GetScheduleByID(ctx context.Context, id int) (*models.Schedule, error) {
	sch, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return sch, nil
}

func (s *ScheduleService) isValidTimeSlot(t string) bool {
	if len(t) >= 5 {
		t = t[:5]
	}

	allowedSlots := map[string]bool{
		"12:00": true,
		"16:00": true,
		"19:30": true,
		"22:00": true,
		"23:50": true,
	}
	return allowedSlots[t]
}
