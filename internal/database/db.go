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

func Connect() {
	once.Do(func() {
		dsn := "host=localhost port=5432 user=postgres password=12345 dbname=movie sslmode=disable"

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

	CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		genre VARCHAR(100) NOT NULL,
		duration INT NOT NULL,
		rating NUMERIC(3, 1) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		phone VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	
	CREATE TABLE IF NOT EXISTS customer_wallets (
		id SERIAL PRIMARY KEY,
		customer_id INT UNIQUE NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
		balance NUMERIC(10, 2) DEFAULT 0.00,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	
	CREATE TABLE IF NOT EXISTS halls (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		total_seats INT NOT NULL
	);

	
	CREATE TABLE IF NOT EXISTS seats (
		id SERIAL PRIMARY KEY,
		hall_id INT NOT NULL REFERENCES halls(id) ON DELETE CASCADE,
		row_num INT NOT NULL,
		seat_num INT NOT NULL,
		UNIQUE(hall_id, row_num, seat_num)
	);

	CREATE TABLE IF NOT EXISTS schedules (
		id SERIAL PRIMARY KEY,
		movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
		hall_id INT REFERENCES halls(id) ON DELETE SET NULL,
		date_time TIMESTAMP NOT NULL,
		price NUMERIC(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tickets (
		id SERIAL PRIMARY KEY,
		schedule_id INT NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
		customer_id INT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
		seat_id INT NOT NULL REFERENCES seats(id) ON DELETE CASCADE,
		status VARCHAR(50) DEFAULT 'booked', -- booked, paid, cancelled
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(schedule_id, seat_id) -- Чтобы одно сиденье на сеанс нельзя было купить дважды
	);

	CREATE TABLE IF NOT EXISTS refunds (
		id SERIAL PRIMARY KEY,
		ticket_id INT UNIQUE NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
		amount NUMERIC(10, 2) NOT NULL,
		refunded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Ошибка при выполнении миграции всех таблиц: %v", err)
	}
	log.Println("Миграции базы данных успешно применены! Все таблицы кинотеатра (включая залы, кошельки и возвраты) готовы к работе.")
}
