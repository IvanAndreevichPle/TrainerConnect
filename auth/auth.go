package auth

import (
	"TrainerConnect/internal/user"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// TokenClaims представляет пользовательские данные, которые могут быть добавлены в токен
type TokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// Authenticate проверяет учетные данные пользователя и возвращает токен
func Authenticate(username, password string, storage *user.Storage, secretKey string) (string, error) {
	// Получаем пользователя из базы данных по имени пользователя
	u, err := storage.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	// Сравниваем введенный пароль с хэшированным паролем из базы данных
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", err
	}

	// Генерируем токен с информацией о пользователе и сроке действия токена
	claims := TokenClaims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Пример: токен действителен 24 часа
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// AuthMiddleware проверяет наличие и валидность токена в заголовке Authorization
func AuthMiddleware(next http.Handler, storage *user.Storage, secretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Токен должен начинаться с "Bearer "
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Проверяем валидность токена и получаем информацию о пользователе
		claims, err := ValidateToken(tokenString, secretKey)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Получаем пользователя из базы данных по имени пользователя
		username := claims.Username
		u, err := storage.GetUserByUsername(username)
		if err != nil {
			http.Error(w, "Error getting user by username", http.StatusInternalServerError)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user", u))

		next.ServeHTTP(w, r)
	})
}

// ValidateToken проверяет валидность токена и возвращает информацию о пользователе, если токен действителен
func ValidateToken(tokenString, secretKey string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}
