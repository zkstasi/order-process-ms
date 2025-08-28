package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Глобальные клиенты
var (
	MongoClient         *mongo.Client
	OrderCollection     *mongo.Collection
	UserCollection      *mongo.Collection
	DeliveryCollection  *mongo.Collection
	WarehouseCollection *mongo.Collection
	RedisClient         *redis.Client
	Ctx                 = context.Background()
)

// InitDB подключает MongoDB и Redis
func InitDB() error {
	// MongoDB
	mongoURI := "mongodb://root:example@localhost:27017/?authSource=admin"
	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("ошибка подключения к MongoDB: %w", err)
	}

	// Проверка соединения с Mongo
	ctx, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("не удалось подключиться к MongoDB: %w", err)
	}

	MongoClient = client
	OrderCollection = client.Database("orderdb").Collection("orders")
	UserCollection = client.Database("orderdb").Collection("users")
	DeliveryCollection = client.Database("orderdb").Collection("deliveries")
	WarehouseCollection = client.Database("orderdb").Collection("warehouses")
	fmt.Println("MongoDB подключена успешно")

	// Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // порт Redis
		Password: "",               // если есть пароль
		DB:       0,                // используем базу 0
	})

	// проверяем соединение
	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("не удалось подключиться к Redis: %w", err)
	}
	fmt.Println("Redis подключен успешно")
	return nil
}

// CloseDB закрывает соединения при завершении работы приложения
func CloseDB() {
	if MongoClient != nil {
		if err := MongoClient.Disconnect(Ctx); err != nil {
			log.Println("Ошибка при отключении MongoDB:", err)
		}
	}
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Println("Ошибка при отключении Redis:", err)
		}
	}
}
