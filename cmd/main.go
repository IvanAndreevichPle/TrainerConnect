package main

import (
	"TrainerConnect/internal/user"
	postgres "TrainerConnect/pkg/postgresql"
	"TrainerConnect/pkg/postgresql/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.ReadConfig("../pkg/postgresql/config/database.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем экземпляр *user.Storage, передавая *sql.DB
	storage := user.NewStorage(db)

	router := setupRouter(storage)
	startServer(router)
}

func setupRouter(storage *user.Storage) *chi.Mux {
	router := chi.NewRouter()

	// Добавляем базовые middleware, такие, как логирование
	router.Use(middleware.Logger)

	// Создаем экземпляр *user.Handler, передавая *user.Storage
	userHandler := user.NewHandler(storage)

	// Регистрируем обработчик в созданном ранее маршрутизаторе
	userHandler.Register(router)

	return router
}

func startServer(router *chi.Mux) {
	server := &http.Server{
		Addr:         "0.0.0.0:1234",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
