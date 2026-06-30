package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"final-project/internal/database"
	"final-project/internal/errs"
	"final-project/internal/models"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{db: database.DB}
}

// 1. ТРАНЗАКЦИЯ: Создание клиента + Автосоздание пустого кошелька
func (r *CustomerRepository) CreateWithWalletTx(ctx context.Context, c *models.Customer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Шаг А: Создаем клиента
	customerQuery := `INSERT INTO customers (nickname, password_hash, phone) VALUES ($1, $2, $3) RETURNING customer_id`
	err = tx.QueryRowContext(ctx, customerQuery, c.Nickname, c.PasswordHash, c.Phone).Scan(&c.ID)
	if err != nil {
		return err // Если никнейм уже занят, база вернет ошибку дубликата UNIQUE
	}

	// Шаг Б: Автоматически создаем ему кошелек с балансом 0
	walletQuery := `INSERT INTO customer_wallet (customer_id, balance) VALUES ($1, 0)`
	_, err = tx.ExecContext(ctx, walletQuery, c.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// 2. Получение клиента по никнейму (нужно для авторизации Login)
func (r *CustomerRepository) GetByNickname(ctx context.Context, nickname string) (*models.Customer, error) {
	var c models.Customer
	query := `SELECT customer_id, nickname, password_hash, phone FROM customers WHERE nickname = $1`
	err := r.db.QueryRowContext(ctx, query, nickname).Scan(&c.ID, &c.Nickname, &c.PasswordHash, &c.Phone)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// 3. Получение профиля по ID вместе с балансом кошелька через JOIN
func (r *CustomerRepository) GetProfileByID(ctx context.Context, id int) (*models.CustomerProfileResponse, error) {
	var p models.CustomerProfileResponse
	query := `
		SELECT c.customer_id, c.nickname, c.phone, w.balance 
		FROM customers c
		JOIN customer_wallet w ON c.customer_id = w.customer_id
		WHERE c.customer_id = $1`
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Nickname, &p.Phone, &p.Balance)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.New("клиент не найден", 404)
	}
	return &p, err
}

// 4. Пополнение кошелька (приход денег)
func (r *CustomerRepository) AddMoney(ctx context.Context, customerID int, amount int) (int64, error) {
	query := `UPDATE customer_wallet SET balance = balance + $1 WHERE customer_id = $2`
	res, err := r.db.ExecContext(ctx, query, amount, customerID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 5. Динамический PATCH для обновления любых полей (телефон, никнейм, хэш пароля)
func (r *CustomerRepository) UpdateFields(ctx context.Context, id int, fields map[string]interface{}) (int64, error) {
	if len(fields) == 0 {
		return 0, nil
	}

	query := "UPDATE customers SET "
	var args []interface{}
	idx := 1

	for key, val := range fields {
		query += key + " = $" + strconv.Itoa(idx) + ", "
		args = append(args, val)
		idx++
	}
	query = query[:len(query)-2] // Отрезаем лишнюю запятую
	query += " WHERE customer_id = $" + strconv.Itoa(idx)
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}


func (r *CustomerRepository) Delete(ctx context.Context, id int) (int64, error) {
	query := `DELETE FROM customers WHERE customer_id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}


func (r *CustomerRepository) GetPaginated(ctx context.Context, limit, offset int) ([]models.CustomerProfileResponse, error) {
	query := `
		SELECT c.customer_id, c.nickname, c.phone, w.balance 
		FROM customers c
		JOIN customer_wallet w ON c.customer_id = w.customer_id
		ORDER BY c.customer_id ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.CustomerProfileResponse
	for rows.Next() {
		var p models.CustomerProfileResponse
		if err := rows.Scan(&p.ID, &p.Nickname, &p.Phone, &p.Balance); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}