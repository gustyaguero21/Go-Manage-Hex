package user

type MysqlRepository interface {
	CreateTable(tableName string) error
}
