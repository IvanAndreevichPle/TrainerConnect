package user

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"math/rand"
	"net/http"
	"strconv"
)

const userURL = "/users/"

type Handler struct {
}

var users = make(map[string]User)

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(router *chi.Mux) {
	users["1"] = User{
		ID:      "1",
		Name:    "Иван",
		Surname: "Иванов",
		Email:   "ivan1@example.com",
	}

	users["2"] = User{
		ID:      "2",
		Name:    "Петр",
		Surname: "Петров",
		Email:   "petr2@example.com",
	}

	users["3"] = User{
		ID:      "3",
		Name:    "Анна",
		Surname: "Сидорова",
		Email:   "anna3@example.com",
	}

	users["4"] = User{
		ID:      "4",
		Name:    "Елена",
		Surname: "Козлова",
		Email:   "elena4@example.com",
	}

	users["5"] = User{
		ID:      "5",
		Name:    "Михаил",
		Surname: "Новиков",
		Email:   "mikhail5@example.com",
	}

	users["6"] = User{
		ID:      "6",
		Name:    "Светлана",
		Surname: "Морозова",
		Email:   "svetlana6@example.com",
	}

	users["7"] = User{
		ID:      "7",
		Name:    "Андрей",
		Surname: "Волков",
		Email:   "andrey7@example.com",
	}

	users["8"] = User{
		ID:      "8",
		Name:    "Татьяна",
		Surname: "Попова",
		Email:   "tatiana8@example.com",
	}

	users["9"] = User{
		ID:      "9",
		Name:    "Ирина",
		Surname: "Семенова",
		Email:   "irina9@example.com",
	}

	users["10"] = User{
		ID:      "10",
		Name:    "Алексей",
		Surname: "Королев",
		Email:   "alexey10@example.com",
	}

	router.Mount(userURL, router)
	router.Get(userURL, h.GetList)
	router.Get(userURL+"{id}", h.GetUserHandler)
	router.Post(userURL, h.CreateNewUserHandler)
	router.Put(userURL+"{id}", h.UpdateUser)
	router.Patch(userURL+"{id}", h.PatchUser)
	router.Delete(userURL+"{id}", h.DeleteUser)
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	// Ищем пользователя по идентификатору в карте users
	user, found := users[id]
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Отправляем данные пользователя в ответ
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) CreateNewUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.ID = strconv.Itoa(rand.Intn(1000000))
	users[user.ID] = user
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "id")
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Найдем пользователя по UUID
	user, found := users[uuid]
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Обновим данные пользователя
	user.Name = updatedUser.Name
	user.Surname = updatedUser.Surname
	user.Email = updatedUser.Email

	// Обновим пользователя в карте users
	users[uuid] = user

	// Отправим ответ с обновленными данными пользователя
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) PatchUser(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "id")
	var patchData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Найдем пользователя по UUID
	user, found := users[uuid]
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Применим частичные изменения к пользователя на основе данных из запроса
	if name, ok := patchData["Name"]; ok {
		user.Name = name
	}
	if surname, ok := patchData["Surname"]; ok {
		user.Surname = surname
	}
	if email, ok := patchData["Email"]; ok {
		user.Email = email
	}

	// Обновим пользователя в карте users
	users[uuid] = user

	// Отправим ответ с обновленными данными пользователя
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "id")

	// Проверим, существует ли пользователь с указанным UUID
	_, found := users[uuid]
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Удалим пользователя из карты users
	delete(users, uuid)

	// Отправим ответ об успешном удалении
	w.Write([]byte("User with UUID " + uuid + " has been successfully deleted"))
}
