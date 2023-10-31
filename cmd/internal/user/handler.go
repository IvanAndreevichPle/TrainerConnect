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
	router.HandleFunc(userURL, h.GetUserByUUID)
	//router.HandleFunc(userURL, h.CreateUser)
	//router.HandleFunc(userURL, h.UpdateUser)
	//router.HandleFunc(userURL, h.DeleteUser)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	// Обработка GET-запроса для получения списка пользователей
	w.Write([]byte("GetList: This is a list of users"))
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) {
	// Обработка GET-запроса для получения пользователя по UUID
	w.Write([]byte("GetUserByUUID: Getting user by UUID"))
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Обработка POST-запроса для создания пользователя
	w.Write([]byte("CreateUser: Creating a new user"))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Обработка PUT-запроса для обновления пользователя
	w.Write([]byte("UpdateUser: Updating user"))
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Обработка DELETE-запроса для удаления пользователя
	w.Write([]byte("DeleteUser: Deleting user"))
}
