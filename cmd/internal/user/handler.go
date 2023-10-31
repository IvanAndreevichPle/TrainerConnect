package user

import (
	"TrainerConnect/cmd/internal/handlers"
	"net/http"
)

const (
	usersURL = "/users"
	userURL  = "/user/:uuid"
)

var _ handlers.Handler = &handler{}

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *http.ServeMux) {
	router.HandleFunc(usersURL, h.GetList)
}

func (h *handler) HandleUser(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get(":uuid")

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
	// Обработка POST-запроса для создания пользователя
	w.Write([]byte("CreateUser: Creating a new user"))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, uuid string) {
	// Обработка PUT-запроса для обновления пользователя
	w.Write([]byte("UpdateUser: Updating user"))
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request, uuid string) {
	// Обработка DELETE-запроса для удаления пользователя
	w.Write([]byte("DeleteUser: Deleting user"))
}
