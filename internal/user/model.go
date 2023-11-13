package user

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	Email     string `json:"email"`
}
