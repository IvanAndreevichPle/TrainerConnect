// user/handler.go

package user

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
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
	router.Get("/users", h.GetUserByUsernameHandler)
	router.Post(userURL, h.CreateNewUserHandler)
	router.Put(userURL+"{id}", h.UpdateUser)
	router.Patch(userURL+"{id}", h.PatchUser)
	router.Delete(userURL+"{id}", h.DeleteUser)
	router.Post("/auth", h.AuthenticateUser)
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

	// Декодирование данных запроса, включая пароль
	var userData struct {
		ID        string `json:"id"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Username  string `json:"username"`
		Role      string `json:"role"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Преобразование в структуру User
	user := User{
		ID:        userData.ID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Username:  userData.Username,
		Role:      userData.Role,
		Email:     userData.Email,
		Password:  userData.Password,
	}

	// Логирование перед созданием пользователя
	log.Printf("Creating new user: %+v", user)

	// Генерация соли
	salt, err := GenerateSalt()
	if err != nil {
		log.Printf("Error generating salt: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Хэширование пароля с использованием соли
	hashedPassword, err := HashPassword(userData.Password, salt)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создание нового пользователя с использованием метода CreateUser
	if err := h.Storage.CreateUser(&user, hashedPassword, salt); err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка ответа с данными созданного пользователя
	json.NewEncoder(w).Encode(user)

	// Логирование после создания пользователя
	log.Printf("User successfully created: %+v", user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Получаем пользователя из базы данных по ID
	updatedUser, err := h.Storage.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Декодируем JSON-тело запроса и обновляем только те поля, которые присутствуют в запросе
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем пользователя в базе данных
	if err := h.Storage.UpdateUser(updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем обновленного пользователя в ответе
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
	if username, ok := patchData["Username"]; ok {
		currentUser.Username = username
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

func (h *Handler) GetUserByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username parameter", http.StatusBadRequest)
		return
	}

	user, _, err := h.Storage.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Error getting user by username", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Отправляем информацию о пользователе в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GenerateSalt генерирует случайную соль
func GenerateSalt() (string, error) {
	saltBytes := make([]byte, 32)
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(saltBytes), nil
}

// HashPassword хэширует пароль с использованием соли
func HashPassword(password, salt string) (string, error) {
	passwordBytes := []byte(password + salt)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// AuthenticateUser аутентифицирует пользователя
func (h *Handler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение пользователя и соли по имени пользователя из базы данных
	existingUser, salt, err := h.Storage.GetUserByUsername(authData.Username)
	if err != nil {
		log.Printf("Error getting user by username: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Сравнение хэша пароля с предоставленным паролем и солью
	err = ComparePasswords(existingUser.Password, authData.Password, salt)
	if err != nil {
		log.Printf("Invalid password for user %s", authData.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// TODO: Создание токена аутентификации (JWT, например) и отправка его в ответе
	// Пример: отправка успешного ответа
	w.Write([]byte("Authentication successful"))
}

// ComparePasswords сравнивает хэш пароля с предоставленным паролем и солью
func ComparePasswords(hashedPassword, password, salt string) error {
	// Проверка пароля
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
	if err != nil {
		return err
	}

	return nil
}
