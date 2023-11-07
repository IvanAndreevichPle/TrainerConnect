package main

import (
	"TrainerConnect/cmd/internal/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	// Создаем новый HTTP-маршрутизатор с использованием Chi
	router := chi.NewRouter()

	// Добавляем базовые middleware, такие, как логгирование
	router.Use(middleware.Logger)

	// Создаем экземпляр для обработчика для ресурса user
	handler := user.NewHandler()

	// Регистрируем обработчик в созданном ранее маршрутизаторе
	handler.Register(router)

	// Запускаем сервер
	start(router)
}

func start(router http.Handler) {
	server := &http.Server{
		Addr:         ":1234",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
