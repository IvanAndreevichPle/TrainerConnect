package user

type User struct {
	ID      string `json:"id"`
	Name    string `json:"firstname"`
	Surname string `json:"lastname"`
	Role    string `json:"role"`
	Email   string `json:"email"`
}
