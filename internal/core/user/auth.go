package user

type Authorization interface {
	GenerateJWT(username, password string) (string, error)
	ValidateJWT(tokenStr string) (string, error)
}
