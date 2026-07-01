package database

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client


func ConnectRedis(addr string) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", 
		DB:       0,  
	})

	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	log.Println("Redis успешно подключен")
}