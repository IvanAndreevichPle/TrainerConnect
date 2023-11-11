// internal/user/handler.go

package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Handler struct {
	Storage *Storage
}

const userURL = "/users/"

func NewHandler(storage *Storage) *Handler {
	return &Handler{Storage: storage}
}

func (h *Handler) Register(router *chi.Mux) {
	router.Mount(userURL, router)
	router.Get(userURL, h.GetList)
	router.Get(userURL+"{id}", h.GetUserHandler)
	router.Post(userURL, h.CreateNewUserHandler)
	router.Put(userURL+"{id}", h.UpdateUser)
	router.Patch(userURL+"{id}", h.PatchUser)
	router.Delete(userURL+"{id}", h.DeleteUser)
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	users, err := h.Storage.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.Storage.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *Handler) CreateNewUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	if err := h.Storage.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser.ID = strconv.Itoa(userID)

	if err := h.Storage.UpdateUser(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedUser)
}

func (h *Handler) PatchUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var patchData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем текущего пользователя
	currentUser, err := h.Storage.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if currentUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Применяем частичные изменения
	if name, ok := patchData["FirstName"]; ok {
		currentUser.FirstName = name
	}
	if surname, ok := patchData["LastName"]; ok {
		currentUser.LastName = surname
	}
	if email, ok := patchData["Email"]; ok {
		currentUser.Email = email
	}

	// Обновляем пользователя
	if err := h.Storage.UpdateUser(currentUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(currentUser)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.Storage.DeleteUser(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User with ID " + id + " has been successfully deleted"))
}
