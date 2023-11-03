package user

type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Role      string `json:"role"`
	Email     string `json:"email"`
}
