package errs

import (
	"sync"
	"time"
)

type AppError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Timestamp  string `json:"timestamp"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(message string, status int) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: status,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
	}
}

type ErrorTracker struct {
	mu           sync.RWMutex
	errorStats   map[string]int
	criticalChan chan *AppError
}

var Tracker *ErrorTracker

func init() {
	Tracker = &ErrorTracker{
		errorStats:   make(map[string]int),
		criticalChan: make(chan *AppError, 100),
	}

	go Tracker.listenCriticalErrors()
}

func (t *ErrorTracker) Track(err *AppError) {
	t.mu.Lock()
	t.errorStats[err.Message]++
	t.mu.Unlock()

	if err.StatusCode == 500 {
		select {
		case t.criticalChan <- err:
		default:

		}
	}
}

func (t *ErrorTracker) listenCriticalErrors() {
	for err := range t.criticalChan {

		println("[CRITICAL ALERT]", err.Timestamp, "-", err.Message)
	}
}

var (
	ErrBadRequest = New("Некорректный запрос", 400)
	ErrNotFound   = New("Ресурс не найден", 404)
	ErrInternal   = New("Внутренняя ошибка сервера", 500)
	ErrTimeout    = New("Время ожидания запроса истекло", 408)

	ErrUnauthorized       = New("Вы не авторизованы", 401)
	ErrInvalidCredentials = New("Неверный логин или пароль", 401)
	ErrInvalidToken       = New("Невалидный или истекший токен", 401)
	ErrNicknameExists     = New("Пользователь с таким никнеймом уже зарегистрирован", 400)

	ErrUserNotFound      = New("Пользователь не найден", 404)
	ErrMovieNotFound     = New("Фильм не найден", 404)
	ErrScheduleNotFound  = New("Сеанс не найден", 404)
	ErrSeatTaken         = New("Это место в зале уже занято на данный сеанс", 400)
	ErrInsufficientFunds = New("Недостаточно средств на балансе кошелька", 400)
	ErrRefundTooLate     = New("Возврат невозможен: до начала сеанса осталось менее 2 часов", 400)
)
