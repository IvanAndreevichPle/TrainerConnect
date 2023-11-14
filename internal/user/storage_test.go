package user_test

import (
	"TrainerConnect/internal/user"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	// Инициализируем моковую БД
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		panic("Error creating mock database: " + err.Error())
	}

	// Запускаем все тесты пока работает моковая БД
	exitCode := m.Run()

	// Закрываем БД после выполнения тестов с БД
	db.Close()

	// Exit with the test result
	os.Exit(exitCode)
}

func TestAllHandlers(t *testing.T) {
	TestCreateNewUserHandler(t)
	TestGetUserHandler(t)
	//TestUpdateUserHandler(t)
	//TestDeleteUserHandler(t)

}
func TestGetUserHandler(t *testing.T) {
	// Set up expectations for the mock database
	mock.ExpectQuery("SELECT first_name, last_name, role, email, username FROM users WHERE user_id = ?").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "role", "email", "username"}).
			AddRow("1", "John", "Doe", "client", "john.doe@example.com", "johndoe"))

	// Create a chi router
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Create a request for the GetUserHandler endpoint
	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()
	fmt.Println(rr)
	// Serve the request to the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var responseUser user.User
	err = json.NewDecoder(rr.Body).Decode(&responseUser)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Check that the response matches the expected user data
	expectedUser := user.User{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "user",
		Email:     "john.doe@example.com",
		Username:  "johndoe",
	}

	assert.Equal(t, expectedUser, responseUser)

	// Verify that the expected SQL query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestCreateNewUserHandler(t *testing.T) {
	// Set up expectations for the mock database
	mock.ExpectExec("INSERT INTO users").
		WithArgs("John", "Doe", "johndoe", sqlmock.AnyArg(), "client", "john.doe@example.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Создаем роутер
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Данные запроса в БД
	requestData := `{
			"firstname": "John",
			"lastname": "Doe",
			"username": "johndoe",
			"role": "client",
			"email": "john.doe@example.com", 
			"password": "sec123"
		}`

	// Формируем POST запрос в тестовую БД
	req, err := http.NewRequest("POST", "/users/", strings.NewReader(requestData))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Задаем заголовок запрос
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	rr := httptest.NewRecorder()

	// Сохраняем запрос в роутер
	router.ServeHTTP(rr, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetListHandler(t *testing.T) {
	// Set up expectations for the mock database
	mock.ExpectQuery("SELECT first_name, last_name, role, email, username FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"first_name", "last_name", "client", "email", "username"}).
			AddRow("John", "Doe", "client", "john.doe@example.com", "johndoe"))
	//AddRow("Jane", "Doe", "client", "jane.doe@example.com", "janedoe"))

	// Create a chi router
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Create a request for the GetListHandler endpoint
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request to the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var responseUsers []user.User
	err = json.NewDecoder(rr.Body).Decode(&responseUsers)
	if err != nil {
		fmt.Println(rr)
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Check that the response matches the expected user data
	expectedUsers := []user.User{
		{FirstName: "John", LastName: "Doe", Role: "user", Email: "john.doe@example.com", Username: "johndoe"},
	}

	assert.Equal(t, expectedUsers, responseUsers)

	// Verify that the expected SQL query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestUpdateUserHandler(t *testing.T) {
	// Set up expectations for the mock database
	mock.ExpectQuery("SELECT user_id, first_name, last_name, role, email, username FROM users WHERE user_id = ?").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "role", "email", "username"}).
			AddRow("1", "John", "Doe", "user", "john.doe@example.com", "johndoe"))

	// Set up expectations for the mock database
	mock.ExpectExec("UPDATE users SET first_name=\\$1, last_name=\\$2, role=\\$3, email=\\$4, username=\\$5 WHERE user_id=\\$6").
		WithArgs("Jane", "Doe", "admin", "jane.doe@example.com", "janedoe", "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a chi router
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Sample request data for updating user
	requestData := `{"firstname": "Jane", "lastname": "Doe", "username": "janedoe", "role": "admin", "email": "jane.doe@example.com"}`

	// Create a request with the sample data
	req, err := http.NewRequest("PUT", "/users/1", strings.NewReader(requestData))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request to the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify that the expected SQL queries were executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestPatchUserHandler(t *testing.T) {
	// Set up expectations for the mock database
	mock.ExpectQuery("SELECT user_id, first_name, last_name, role, email, username FROM users WHERE user_id = ?").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "role", "email", "username"}).
			AddRow("1", "John", "Doe", "user", "john.doe@example.com", "johndoe"))

	// Set up expectations for the mock database
	mock.ExpectExec("UPDATE users SET first_name=\\$1, last_name=\\$2, role=\\$3, email=\\$4, username=\\$5 WHERE user_id=\\$6").
		WithArgs("Jane", "Doe", "admin", "jane.doe@example.com", "janedoe", "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a chi router
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Sample request data for patching user
	requestData := `{"firstname": "Jane", "role": "admin"}`

	// Create a request with the sample data
	req, err := http.NewRequest("PATCH", "/users/1", strings.NewReader(requestData))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request to the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify that the expected SQL queries were executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	// Set up expectations for the mock database
	mock.ExpectExec("DELETE FROM users WHERE user_id = ?").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))

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

	// Verify that the expected SQL queries were executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetAllUsersHandler(t *testing.T) {
	// Set up expectations for the mock database
	mock.ExpectQuery("SELECT user_id, first_name, last_name, role, email, username FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "role", "email", "username"}).
			AddRow("1", "John", "Doe", "user", "john.doe@example.com", "johndoe").
			AddRow("2", "Jane", "Doe", "admin", "jane.doe@example.com", "janedoe"))

	// Create a chi router
	router := chi.NewRouter()
	handler := user.NewHandler(user.NewStorage(db))
	handler.Register(router)

	// Create a request for the GetAllUsersHandler endpoint
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request to the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var responseUsers []user.User
	err = json.NewDecoder(rr.Body).Decode(&responseUsers)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Check that the response matches the expected user data
	expectedUsers := []user.User{
		{ID: "1", FirstName: "John", LastName: "Doe", Role: "user", Email: "john.doe@example.com", Username: "johndoe"},
		{ID: "2", FirstName: "Jane", LastName: "Doe", Role: "admin", Email: "jane.doe@example.com", Username: "janedoe"},
	}

	assert.Equal(t, expectedUsers, responseUsers)

	// Verify that the expected SQL query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
