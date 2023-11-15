package user_test

import (
	"TrainerConnect/internal/user"
	postgres "TrainerConnect/pkg/postgresql"
	"TrainerConnect/pkg/postgresql/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/DATA-DOG/go-sqlmock"
)

var db *sql.DB
var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	// Инициализируем тестовую БД
	cfgTest, err := config.ReadConfig("../../pkg/postgresql/config/database_test.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err = postgres.NewDB(cfgTest)
	if err != nil {
		log.Fatal(err)
	}
	// Закрываем БД после выполнения тестов с БД
	defer db.Close()

	// Запускаем все тесты пока работает моковая БД
	exitCode := m.Run()

	// Exit with the test result
	os.Exit(exitCode)
}

//func TestAllHandlers(t *testing.T) {
//	TestCreateNewUserHandler(t)
//	TestGetUserHandler(t)
//	TestUpdateUserHandler(t)
//	//TestDeleteUserHandler(t)
//}

func TestCreateNewUserHandler(t *testing.T) {
	// Создаем роутер
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Данные запроса в БД
	requestData1 := `{
		"id": "1",
		"firstname": "John",
		"lastname": "Doe",
		"username": "johndoe",
		"role": "client",
		"email": "john.doe@example.com", 
		"password": "sec123"
	}`

	requestData2 := `{
		"id": "2",
		"firstname": "Joahn",
		"lastname": "Doe",
		"username": "joahndoe",
		"role": "client",
		"email": "joahn.doe@example.com", 
		"password": "set123"
	}`

	// Формируем POST запрос в тестовую БД для первого пользователя
	req1, err := http.NewRequest("POST", "/users/", strings.NewReader(requestData1))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Задаем заголовок запроса
	req1.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	rr1 := httptest.NewRecorder()

	// Сохраняем запрос в роутер
	router.ServeHTTP(rr1, req1)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, rr1.Code)

	// Формируем POST запрос в тестовую БД для второго пользователя
	req2, err := http.NewRequest("POST", "/users/", strings.NewReader(requestData2))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Задаем заголовок запроса
	req2.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	rr2 := httptest.NewRecorder()

	// Сохраняем запрос в роутер
	router.ServeHTTP(rr2, req2)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, rr2.Code)
}

func TestGetUserHandler(t *testing.T) {
	// Создаем роутер
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Формируем GET запрос в тестовую БД
	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Записываем ответ
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var responseUser user.User
	err = json.NewDecoder(rr.Body).Decode(&responseUser)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Проверяем, что ответ соответствует ожиданию
	expectedUsers := user.User{
		ID: "1", FirstName: "John", LastName: "Doe", Role: "client", Email: "john.doe@example.com", Username: "johndoe"}
	assert.Equal(t, expectedUsers, responseUser)
}

func TestGetListHandler(t *testing.T) {
	// Создаем роутер
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Формируем GET запрос в тестовую БД
	req, err := http.NewRequest("GET", "/users?username=johndoe", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Записываем ответ
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var responseUser user.User
	err = json.NewDecoder(rr.Body).Decode(&responseUser)
	if err != nil {
		fmt.Println(rr)
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Проверяем, что ответ соответствует ожиданию
	expectedUser := user.User{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "client",
		Email:     "john.doe@example.com",
		Username:  "johndoe",
	}
	compareUser := user.User{
		ID:        responseUser.ID,
		FirstName: responseUser.FirstName,
		LastName:  responseUser.LastName,
		Role:      responseUser.Role,
		Email:     responseUser.Email,
		Username:  responseUser.Username,
	}
	assert.Equal(t, expectedUser, compareUser)
}

func TestUpdateUserHandler(t *testing.T) {
	// Создаем роутер
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Данные запроса в БД
	//user, _, err := h.Storage.GetUserByUsername(username)
	requestData := `{
		"firstname": "John",
		"lastname": "Doe",
		"role": "client",
		"email": "john.doe@example.com"
	}`

	// Формируем PUT запрос в тестовую БД с явным указанием значения username
	req, err := http.NewRequest("PUT", "/users/1", strings.NewReader(requestData))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()
	fmt.Println(rr)
	// Serve the request to the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAllUsersHandler(t *testing.T) {
	// Создаем роутер
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Формируем GET запрос в тестовую БД
	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Записываем ответ
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var responseUser []user.User
	err = json.NewDecoder(rr.Body).Decode(&responseUser)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Проверяем, что ответ соответствует ожиданию
	expectedUsers := []user.User{
		{ID: "2", FirstName: "Joahn", LastName: "Doe", Role: "client", Email: "joahn.doe@example.com", Username: "joahndoe"},
		{ID: "1", FirstName: "John", LastName: "Doe", Role: "client", Email: "john.doe@example.com", Username: "johndoe"},
	}
	assert.Equal(t, expectedUsers, responseUser)
}

func TestDeleteUserHandler(t *testing.T) {
	// Create a chi router
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Create a request for the DeleteUserHandler endpoint
	req, err := http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request to the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

//func TestPatchUserHandler(t *testing.T) {
//	// Set up expectations for the mock database
//	mock.ExpectQuery("SELECT user_id, first_name, last_name, role, email, username FROM users WHERE user_id = ?").
//		WithArgs("1").
//		WillReturnRows(sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "role", "email", "username"}).
//			AddRow("1", "John", "Doe", "user", "john.doe@example.com", "johndoe"))
//
//	// Set up expectations for the mock database
//	mock.ExpectExec("UPDATE users SET first_name=\\$1, last_name=\\$2, role=\\$3, email=\\$4, username=\\$5 WHERE user_id=\\$6").
//		WithArgs("Jane", "Doe", "admin", "jane.doe@example.com", "janedoe", "1").
//		WillReturnResult(sqlmock.NewResult(1, 1))
//
//	// Create a chi router
//	router := chi.NewRouter()
//	handler := user.NewHandler(user.NewStorage(db))
//	handler.Register(router)
//
//	// Sample request data for patching user
//	requestData := `{"firstname": "Joahne", "role": "trainer"}`
//
//	// Create a request with the sample data
//	req, err := http.NewRequest("PATCH", "/users/1", strings.NewReader(requestData))
//	if err != nil {
//		t.Fatalf("Error creating request: %v", err)
//	}
//
//	// Set the content type header
//	req.Header.Set("Content-Type", "application/json")
//
//	// Create a response recorder to capture the response
//	rr := httptest.NewRecorder()
//
//	// Serve the request to the router
//	router.ServeHTTP(rr, req)
//
//	// Check the response status code
//	assert.Equal(t, http.StatusOK, rr.Code)
//
//	// Verify that the expected SQL queries were executed
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("Unfulfilled expectations: %s", err)
//	}
//}
//
