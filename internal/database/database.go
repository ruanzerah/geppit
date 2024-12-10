package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ruanzerah/geppit/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	GetQueries() *repository.Queries
	Close() error
}

type service struct {
	Queries *repository.Queries
	Pool    *pgxpool.Pool
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Failed to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	dbInstance = &service{
		Queries: repository.New(pool),
		Pool:    pool,
	}
	return dbInstance
}

func (s *service) GetQueries() *repository.Queries {
	return s.Queries
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.Pool.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := s.Pool.Stat()
	stats["total_connections"] = strconv.Itoa(int(dbStats.TotalConns()))
	stats["idle_connections"] = strconv.Itoa(int(dbStats.IdleConns()))
	stats["used_connections"] = strconv.Itoa(int(dbStats.AcquiredConns()))
	stats["max_connections"] = strconv.Itoa(int(dbStats.MaxConns()))
	stats["canceled"] = strconv.Itoa(int(dbStats.CanceledAcquireCount()))
	stats["acquire_count"] = strconv.Itoa(int(dbStats.AcquireCount()))
	stats["acquire_duration"] = dbStats.AcquireDuration().String()

	if dbStats.TotalConns() > 50 {
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.CanceledAcquireCount() > 100 {
		stats["message"] = "High number of canceled acquire attempts, check for contention."
	}

	return stats
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	s.Pool.Close()
	return nil
}
