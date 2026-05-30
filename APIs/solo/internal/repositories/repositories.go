package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repositories struct {
	DB *sql.DB
}

func New() (*Repositories, error) {
	db, err := openDBFromEnv()
	if err != nil {
		return nil, err
	}

	return &Repositories{DB: db}, nil
}

func openDBFromEnv() (*sql.DB, error) {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		host := strings.TrimSpace(os.Getenv("DB_HOST"))
		port := strings.TrimSpace(os.Getenv("DB_PORT"))
		user := strings.TrimSpace(os.Getenv("DB_USER"))
		password := os.Getenv("DB_PASSWORD")
		database := strings.TrimSpace(os.Getenv("DB_NAME"))
		sslmode := strings.TrimSpace(os.Getenv("DB_SSLMODE"))

		if sslmode == "" {
			sslmode = "disable"
		}

		if host == "" || port == "" || user == "" || database == "" {
			return nil, fmt.Errorf("database configuration is incomplete")
		}

		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, database, sslmode)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
