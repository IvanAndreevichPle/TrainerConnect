package user

import (
	"database/sql"
	"log"
)

type Storage struct {
	*sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

// CreateUser создает нового пользователя в базе данных
func (s *Storage) CreateUser(user *User, password, salt string) error {
	_, err := s.DB.Exec("INSERT INTO users (first_name, last_name, username, password, role, email, salt) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.FirstName, user.LastName, user.Username, password, user.Role, user.Email, salt)

	return err
}

func (s *Storage) GetUserByID(userID int) (*User, error) {
	// Реализация получения пользователя из базы данных по ID
	row := s.DB.QueryRow("SELECT user_id, first_name, last_name, role, email, username FROM users WHERE user_id = $1", userID)
	user := &User{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Storage) UpdateUser(user *User) error {
	// Реализация обновления данных пользователя в базе данных
	_, err := s.DB.Exec("UPDATE users SET first_name=$1, last_name=$2, role=$3, email=$4, username=$5 WHERE user_id=$6",
		user.FirstName, user.LastName, user.Role, user.Email, user.Username, user.ID)
	return err
}

func (s *Storage) DeleteUser(userID int) error {
	// Реализация удаления пользователя из базы данных
	_, err := s.DB.Exec("DELETE FROM users WHERE user_id = $1", userID)
	return err
}

func (s *Storage) GetAllUsers() ([]User, error) {
	// Реализация получения списка всех пользователей из базы данных
	rows, err := s.DB.Query("SELECT user_id, first_name, last_name, role, email, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.Email, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUserByUsername возвращает пользователя и соль по имени пользователя из базы данных
func (s *Storage) GetUserByUsername(username string) (*User, string, error) {
	query := "SELECT user_id, username, password, salt, role, first_name, last_name, email FROM users WHERE username = $1"
	row := s.DB.QueryRow(query, username)

	var u User

	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Salt, &u.Role, &u.FirstName, &u.LastName, &u.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil // Пользователь не найден, возвращаем nil и пустую соль
		}
		log.Printf("Error getting user by username: %v", err)
		return nil, "", err
	}

	return &u, u.Salt, nil
}
