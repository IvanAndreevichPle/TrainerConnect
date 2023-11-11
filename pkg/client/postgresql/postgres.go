package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// Config представляет конфигурацию для подключения к базе данных PostgreSQL
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// NewDB создает новое подключение к базе данных PostgreSQL
func NewDB(cfg Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
