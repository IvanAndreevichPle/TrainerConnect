package user

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Role      string `json:"role"`
	Email     string `json:"email"`
}
