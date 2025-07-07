package user

type Authorization interface {
	GenerateJWT(username string) (string, error)
	ValidateJWT(tokenStr string) (string, error)
}
