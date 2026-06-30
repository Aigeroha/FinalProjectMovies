package services

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"final-project/internal/errs"
	"final-project/internal/models"
	"final-project/internal/repository"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(r *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: r}
}

func (s *CustomerService) Register(ctx context.Context, input models.RegisterInput) error {
	if input.Nickname == "" || input.Password == "" {
		return errs.New("nickname и password не могут быть пустыми", 400)
	}

	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errs.ErrInternal
	}

	customer := &models.Customer{
		Nickname:     input.Nickname,
		PasswordHash: string(hashedPassword),
		Phone:        input.Phone,
	}

	
	err = s.repo.CreateWithWalletTx(ctx, customer)
	if err != nil {
		return err
	}

	return nil
}


func (s *CustomerService) Login(ctx context.Context, input models.LoginInput) (*models.Customer, error) {
	customer, err := s.repo.GetByNickname(ctx, input.Nickname)
	if err != nil {
		return nil, errs.New("неверный никнейм или пароль", 401)
	}

	
	err = bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, errs.New("неверный никнейм или пароль", 401)
	}

	return customer, nil
}


func (s *CustomerService) DeleteCustomer(ctx context.Context, id int) error {
	rowsAffected, err := s.repo.Delete(ctx, id)
	if err != nil {
		return errs.ErrInternal
	}
	if rowsAffected == 0 {
		return errs.New("клиент не найден", 404)
	}
	return nil
}


func (s *CustomerService) PatchCustomer(ctx context.Context, id int, input map[string]interface{}) error {

	if passwordPlain, ok := input["password"].(string); ok && passwordPlain != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(passwordPlain), bcrypt.DefaultCost)
		if err != nil {
			return errs.ErrInternal
		}
		delete(input, "password") 
		input["password_hash"] = string(hashed)
	}

	rowsAffected, err := s.repo.UpdateFields(ctx, id, input)
	if err != nil {
		return errs.ErrInternal
	}
	if rowsAffected == 0 {
		return errs.New("клиент не найден", 404)
	}
	return nil
}


func (s *CustomerService) TopUpWallet(ctx context.Context, customerID int, amount int) error {
	if amount <= 0 {
		return errs.New("сумма пополнения должна быть больше нуля", 400)
	}

	rowsAffected, err := s.repo.AddMoney(ctx, customerID, amount)
	if err != nil {
		return errs.ErrInternal
	}
	if rowsAffected == 0 {
		return errs.New("кошелек для данного клиента не найден", 404)
	}
	return nil
}


func (s *CustomerService) GetCustomerProfile(ctx context.Context, id int) (*models.CustomerProfileResponse, error) {
	profile, err := s.repo.GetProfileByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return profile, nil
}


func (s *CustomerService) GetAllCustomersPaginated(ctx context.Context, page int) ([]models.CustomerProfileResponse, error) {
	limit := 20
	offset := (page - 1) * limit

	customers, err := s.repo.GetPaginated(ctx, limit, offset)
	if err != nil {
		return nil, errs.ErrInternal
	}
	return customers, nil
}