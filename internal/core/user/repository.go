package user

type MysqlRepository interface {
	CreateTable(tableName string) error
	GetByName(name string) (User, error)
	NewUser(user User) error
	DeleteUser(name string) error
	UpdateUser(name string, user User) error
	ChangePwd(newPwd, name string) error
}
