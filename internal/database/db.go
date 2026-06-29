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

		log.Println("PostgreSQL база данных 'movie' успешно подключена (Инициализация Once)")
	})
}
