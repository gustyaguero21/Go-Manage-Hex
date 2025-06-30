package user

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
