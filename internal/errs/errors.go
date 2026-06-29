package errs

type AppError struct {
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func New(message string, status int) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: status,
	}

}

var (
	ErrBadRequest = New("Bad request", 400)
	ErrNotFound   = New("Resource not found", 404)
	ErrInternal   = New("Internal server error", 500)

	ErrUserNotFound       = New("User not found", 404)
	ErrEmailExists        = New("Email already exists", 400)
	ErrTimeout            = New("REquest timeout", 408)
	ErrUnauthorized       = New("Unauthorized", 401)
	ErrInvalidCredentials = New("Invalid credentials", 401)
	ErrInvalidToken       = New("Invalid token", 401)
)