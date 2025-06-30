package user

type MysqlRepository interface {
	CreateTable(tableName string) error
	GetByUsername(username string) (User, error)
	CheckExists(username string) bool
	NewUser(user User) error
	DeleteUser(username string) error
	UpdateUser(username string, user User) error
	ChangePwd(newPwd, username string) error
}
