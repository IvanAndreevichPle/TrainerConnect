package user

import (
	"TrainerConnect/cmd/internal/handlers"
	"encoding/json"
	"net/http"
)

const (
	usersURL = "/users/"
)

var _ handlers.Handler = &handler{}

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *http.ServeMux) {
	router.HandleFunc(usersURL, h.HandleUser)
}

func (h *handler) HandleUser(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")

	switch r.Method {
	case http.MethodGet:
		h.GetUserByUUID(w, r, uuid)
	case http.MethodPost:
		h.CreateUser(w, r, uuid)
	case http.MethodPut:
		h.UpdateUser(w, r, uuid)
	case http.MethodDelete:
		h.DeleteUser(w, r, uuid)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	// Обработка GET-запроса для получения списка пользователей
	w.Write([]byte("GetList: This is a list of users"))
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request, uuid string) {
	// Обработка GET-запроса для получения пользователя по UUID
	w.Write([]byte("GetUserByUUID: Getting user by UUID"))
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, uuid string) {
	if uuid != "" {
		http.Error(w, "UUID should not be provided for user creation", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверки на валидность входных данных и создание пользователя

	// Отправка ответа с созданным пользователем
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, uuid string) {
	if uuid == "" {
		http.Error(w, "UUID is required for user update", http.StatusBadRequest)
		return
	}

	// Обработка PUT-запроса для обновления пользователя
	w.Write([]byte("UpdateUser: Updating user"))
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request, uuid string) {
	if uuid == "" {
		http.Error(w, "UUID is required for user deletion", http.StatusBadRequest)
		return
	}

	// Обработка DELETE-запроса для удаления пользователя
	w.Write([]byte("DeleteUser: Deleting user"))
}
