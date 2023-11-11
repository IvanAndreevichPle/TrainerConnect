// internal/user/storage.go

package user

import (
	"database/sql"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) CreateUser(user *User) error {
	// Реализация создания пользователя в базе данных
	_, err := s.DB.Exec("INSERT INTO users (first_name, last_name, role, email) VALUES ($1, $2, $3, $4)",
		user.FirstName, user.LastName, user.Role, user.Email)
	return err
}

func (s *Storage) GetUserByID(userID int) (*User, error) {
	// Реализация получения пользователя из базы данных по ID
	row := s.DB.QueryRow("SELECT user_id, first_name, last_name, role, email FROM users WHERE user_id = $1", userID)
	user := &User{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Storage) UpdateUser(user *User) error {
	// Реализация обновления данных пользователя в базе данных
	_, err := s.DB.Exec("UPDATE users SET first_name=$1, last_name=$2, role=$3, email=$4 WHERE user_id=$5",
		user.FirstName, user.LastName, user.Role, user.Email, user.ID)
	return err
}

func (s *Storage) DeleteUser(userID int) error {
	// Реализация удаления пользователя из базы данных
	_, err := s.DB.Exec("DELETE FROM users WHERE user_id = $1", userID)
	return err
}

func (s *Storage) GetAllUsers() ([]User, error) {
	// Реализация получения списка всех пользователей из базы данных
	rows, err := s.DB.Query("SELECT user_id, first_name, last_name, role, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
