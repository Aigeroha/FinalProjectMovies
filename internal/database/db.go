package database

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"final-project/internal/errs"

	_ "github.com/lib/pq"
)

var (
	DB   *sql.DB
	once sync.Once
)

func Connect(dsn string) {
	once.Do(func() {

		var err error
		DB, err = sql.Open("postgres", dsn)
		if err != nil {
			appErr := errs.New("Ошибка конфигурации БД: "+err.Error(), 500)
			errs.Tracker.Track(appErr)
			log.Fatal(appErr.Error())
		}

		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(25)
		DB.SetConnMaxLifetime(5 * time.Minute)

		if err = DB.Ping(); err != nil {
			appErr := errs.New("БД не подключена (Ping failed): "+err.Error(), 500)
			errs.Tracker.Track(appErr)
			log.Fatal(appErr.Error())
		}

		log.Println("PostgreSQL база данных успешно подключена")

		runMigrations()
	})
}

func runMigrations() {
	query := `

	CREATE TABLE IF NOT EXISTS customers (
		customer_id SERIAL PRIMARY KEY,
		nickname VARCHAR(100) UNIQUE,
		password_hash VARCHAR(255),
		phone VARCHAR(25)
	);

	CREATE TABLE IF NOT EXISTS movies (
		movie_id SERIAL PRIMARY KEY,
		title VARCHAR(100),
		duration INT,
		genre VARCHAR(50),
		rating DECIMAL(2,1)
	);

	CREATE TABLE IF NOT EXISTS halls (
		hall_id SERIAL PRIMARY KEY,
		hall_number INT UNIQUE,
		capacity INT
	);

	CREATE TABLE IF NOT EXISTS seats (
		seat_id SERIAL PRIMARY KEY,
		hall_id INT REFERENCES halls(hall_id) ON DELETE CASCADE,
		row_number INT CHECK (row_number BETWEEN 1 AND 3),
		seat_number INT CHECK (seat_number BETWEEN 1 AND 10)
	);

	CREATE TABLE IF NOT EXISTS customer_wallet (
		account_id SERIAL PRIMARY KEY,
		customer_id INT UNIQUE REFERENCES customers(customer_id) ON DELETE CASCADE,
		balance INT CHECK (balance >= 0)
	);

	CREATE TABLE IF NOT EXISTS schedules (
		schedule_id SERIAL PRIMARY KEY,
		movie_id INT REFERENCES movies(movie_id) ON DELETE CASCADE,
		session_date DATE,
		session_time TIME,
		hall_id INT REFERENCES halls(hall_id) ON DELETE SET NULL,
		adult_price INT,
		student_price INT,
		child_price INT
	);

	CREATE TABLE IF NOT EXISTS tickets (
		ticket_id SERIAL PRIMARY KEY,
		schedule_id INT REFERENCES schedules(schedule_id) ON DELETE CASCADE,
		seat_id INT REFERENCES seats(seat_id) ON DELETE CASCADE,
		customer_id INT REFERENCES customers(customer_id) ON DELETE CASCADE,
		ticket_type VARCHAR(20) NOT NULL CHECK (ticket_type IN ('Взрослый', 'Студенческий', 'Детский')),
		status VARCHAR(20) NOT NULL CHECK (status IN ('Куплено', 'Отмена'))
	);

	CREATE TABLE IF NOT EXISTS refunds (
		refund_id SERIAL PRIMARY KEY,
		ticket_id INT UNIQUE REFERENCES tickets(ticket_id) ON DELETE CASCADE,
		refund_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		amount DECIMAL(10,2)
	);

	
	CREATE OR REPLACE VIEW view_customer_wallets_readable AS
	SELECT 
		w.account_id,
		c.nickname AS customer_nickname,
		w.balance
	FROM customer_wallet w
	JOIN customers c ON w.customer_id = c.customer_id;

	CREATE OR REPLACE VIEW view_readable_schedules AS
	SELECT 
		s.schedule_id,
		m.title AS movie_title, 
		s.session_date,
		s.session_time,
		s.hall_id,
		s.adult_price,
		s.student_price,
		s.child_price
	FROM schedules s
	JOIN movies m ON s.movie_id = m.movie_id;
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Ошибка при выполнении миграции всех таблиц и view: %v", err)
	}
	log.Println("Миграции базы данных успешно применены! Все таблицы и view кинотеатраготовы к работе.")
}
