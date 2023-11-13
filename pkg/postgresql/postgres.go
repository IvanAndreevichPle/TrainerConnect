package postgres

import (
	"TrainerConnect/pkg/postgresql/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// NewDB создает новое подключение к базе данных PostgreSQL
func NewDB(cfg config.DBConfig) (*sql.DB, error) {
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
