package user

import (
	"database/sql"
	"go-manage-hex/cmd/config"
	mysqlrepo "go-manage-hex/internal/core/user"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTable(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	repo := NewUserMysql(db)

	test := []struct {
		Name        string
		ExpectedErr error
		MockFunc    func()
	}{
		{
			Name:        "CreateTable_Success",
			ExpectedErr: nil,
			MockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(config.CreateTableTest)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:        "CreateTable_Err",
			ExpectedErr: err,
			MockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(config.CreateTableTest)).
					WillReturnError(err)
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockFunc()

			err := repo.CreateTable("table_name")
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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
				mock.ExpectQuery(config.GetByNameTest).
					WithArgs("John").WillReturnRows(rows)
			},
		},
		{
			Name:         "GetByName_Err",
			SearchName:   "John",
			ExpectedUser: mysqlrepo.User{},
			ExpectedErr:  sql.ErrNoRows,
			MockFunc: func(u mysqlrepo.User) {
				mock.ExpectQuery(config.GetByNameTest).
					WithArgs("John").WillReturnError(err)
			},
		},
		{
			Name:         "GetByName_NoRows",
			SearchName:   "John",
			ExpectedUser: mysqlrepo.User{},
			ExpectedErr:  sql.ErrNoRows,
			MockFunc: func(u mysqlrepo.User) {
				mock.ExpectQuery(config.GetByNameTest).
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

func TestNewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	repo := NewUserMysql(db)

	test := []struct {
		Name        string
		NewUser     mysqlrepo.User
		ExpectedErr error
		MockFunc    func()
	}{
		{
			Name: "NeUser_Success",
			NewUser: mysqlrepo.User{
				ID:       "1",
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234.",
			},
			ExpectedErr: nil,
			MockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(config.NewUserTest)).
					WithArgs("1", "John", "Doe", "johndoe", "johndoe@example.com", "Password1234.").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name: "NeUser_Err",
			NewUser: mysqlrepo.User{
				ID:       "1",
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234.",
			},
			ExpectedErr: err,
			MockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(config.NewUserTest)).WillReturnError(err)
			},
		},
	}
	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockFunc()

			err := repo.NewUser(tt.NewUser)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	repo := NewUserMysql(db)

	test := []struct {
		Name        string
		DeleteUser  string
		ExpectedErr error
		MockFunc    func()
	}{
		{
			Name:        "DeleteUser_Success",
			DeleteUser:  "John",
			ExpectedErr: nil,
			MockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(config.DeleteUserTest)).
					WithArgs("John").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:        "DeleteUser_Err",
			DeleteUser:  "John",
			ExpectedErr: err,
			MockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(config.DeleteUserTest)).
					WithArgs("John").
					WillReturnError(err)
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockFunc()

			err := repo.DeleteUser(tt.DeleteUser)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	repo := NewUserMysql(db)

	test := []struct {
		Name        string
		UpdateUser  mysqlrepo.User
		Update      string
		ExpectedErr error
		MockFunc    func()
	}{
		{
			Name: "UpdateUser_Success",
			UpdateUser: mysqlrepo.User{
				ID:       "1",
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234.",
			},
			Update:      "John",
			ExpectedErr: nil,
			MockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(config.UpdateUserTest)).
					WithArgs("John", "Doe", "johndoe@example.com", "John").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:        "UpdateUser_Err",
			UpdateUser:  mysqlrepo.User{},
			Update:      "John",
			ExpectedErr: err,
			MockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(config.UpdateUserTest)).
					WillReturnError(err)
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockFunc()

			err := repo.UpdateUser(tt.Update, tt.UpdateUser)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChangePwd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	repo := NewUserMysql(db)

	test := []struct {
		Name               string
		ChangePwdUser_name string
		NewPwd             string
		ExpectedErr        error
		MockFunc           func()
	}{
		{
			Name:               "ChangePwd_Success",
			ChangePwdUser_name: "John",
			NewPwd:             "NewPassword",
			ExpectedErr:        nil,
			MockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(config.ChangePwdTest)).
					WithArgs("NewPassword", "John").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:               "ChangePwd_Err",
			ChangePwdUser_name: "John",
			NewPwd:             "NewPassword",
			ExpectedErr:        err,
			MockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(config.ChangePwdTest)).
					WillReturnError(err)
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockFunc()

			err := repo.ChangePwd(tt.NewPwd, tt.ChangePwdUser_name)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
