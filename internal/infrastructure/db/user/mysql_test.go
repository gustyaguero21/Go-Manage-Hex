package user

import (
	"database/sql"
	mysqlrepo "go-manage-hex/internal/core/user"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewUserMysql(db)

	tests := []struct {
		Name         string
		SearchName   string
		ExpectedUser mysqlrepo.User
		ExpectedErr  error
		MockFunc     func(u mysqlrepo.User)
	}{
		{
			Name:         "GetByName_Success",
			SearchName:   "John",
			ExpectedUser: mysqlrepo.User{},
			ExpectedErr:  nil,
			MockFunc: func(u mysqlrepo.User) {
				rows := sqlmock.NewRows([]string{"id", "name", "last_name", "username", "email", "password"}).
					AddRow(u.ID, u.Name, u.LastName, u.Username, u.Email, u.Password)
				mock.ExpectQuery("SELECT id,name,last_name,username,email,password FROM  WHERE name = ?").
					WithArgs("John").WillReturnRows(rows)
			},
		},
		{
			Name:         "GetByName_Err",
			SearchName:   "John",
			ExpectedUser: mysqlrepo.User{},
			ExpectedErr:  sql.ErrNoRows,
			MockFunc: func(u mysqlrepo.User) {
				mock.ExpectQuery("SELECT id,name,last_name,username,email,password FROM  WHERE name = ?").
					WithArgs("John").WillReturnError(err)
			},
		},
		{
			Name:         "GetByName_NoRows",
			SearchName:   "John",
			ExpectedUser: mysqlrepo.User{},
			ExpectedErr:  sql.ErrNoRows,
			MockFunc: func(u mysqlrepo.User) {
				mock.ExpectQuery("SELECT id,name,last_name,username,email,password FROM  WHERE name = ?").
					WithArgs("John").WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockFunc(tt.ExpectedUser)

			result, err := repo.GetByName(tt.SearchName)

			if tt.ExpectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedUser.ID, result.ID)
			}
		})
	}
}
