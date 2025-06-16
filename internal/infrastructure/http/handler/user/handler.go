package user

import (
	"go-manage-hex/cmd/config"
	"go-manage-hex/internal/app/user"
	"net/http"
	"strconv"

	entity "go-manage-hex/internal/core/user"

	"github.com/gin-gonic/gin"
	"github.com/gustyaguero21/go-core/pkg/web"
)

type UserHandler struct {
	Service user.Usecases
}

func NewUserHandler(service user.Usecases) *UserHandler {
	return &UserHandler{Service: service}
}

func (uh *UserHandler) SearchUserHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	username := c.Query("username")

	if username == "" {
		web.NewError(c, http.StatusBadRequest, config.InvalidQueryParamsMsg)
		return
	}

	search, searchErr := uh.Service.SearchUser(c, username)
	if searchErr != nil {
		web.NewError(c, http.StatusInternalServerError, searchErr.Error())
		return
	}

	c.JSON(http.StatusOK, userResponse(http.StatusOK, config.UserFoundMsg, search))
}

func (uh *UserHandler) CreateUserHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var user entity.User

	if err := c.ShouldBindJSON(&user); err != nil {
		web.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	if user.Name == "" || user.LastName == "" || user.Username == "" || user.Email == "" || user.Password == "" {
		web.NewError(c, http.StatusBadRequest, config.InvalidBodyMsg)
		return
	}

	created, createdErr := uh.Service.CreateUser(c, user)
	if createdErr != nil {
		web.NewError(c, http.StatusInternalServerError, createdErr.Error())
		return
	}

	c.JSON(http.StatusOK, userResponse(http.StatusCreated, config.UserCreatedMsg, created))
}

func (uh *UserHandler) DeleteUserHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	username := c.Query("username")

	if username == "" {
		web.NewError(c, http.StatusBadRequest, config.InvalidQueryParamsMsg)
		return
	}

	confirmation, err := strconv.ParseBool(c.Query("confirmation"))
	if err != nil {
		web.NewError(c, http.StatusBadRequest, config.InvalidConfirmationMsg)
		return
	}

	if !confirmation {
		c.JSON(http.StatusOK, userResponse(http.StatusOK, config.InvalidConfirmationMsg, nil))
		return
	}

	deleteErr := uh.Service.DeleteUser(c, username)
	if deleteErr != nil {
		web.NewError(c, http.StatusInternalServerError, deleteErr.Error())
		return
	}

	c.JSON(http.StatusOK, userResponse(http.StatusOK, config.UserDeletedMsg, nil))
}

func (uh *UserHandler) UpdateUserHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	username := c.Query("username")
	if username == "" {
		web.NewError(c, http.StatusBadRequest, config.InvalidQueryParamsMsg)
		return
	}

	var user entity.User

	if err := c.ShouldBindJSON(&user); err != nil {
		web.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	if user.Name == "" || user.LastName == "" || user.Email == "" {
		web.NewError(c, http.StatusBadRequest, config.InvalidBodyMsg)
		return
	}

	_, updateErr := uh.Service.UpdateUser(c, username, user)
	if updateErr != nil {
		web.NewError(c, http.StatusInternalServerError, updateErr.Error())
		return
	}

	c.JSON(http.StatusOK, userResponse(http.StatusOK, config.UserUpdatedMsg, nil))

}

func (uh *UserHandler) ChangePwdHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	username := c.Query("username")
	newPwd := c.Query("newPwd")

	if username == "" || newPwd == "" {
		web.NewError(c, http.StatusBadRequest, config.InvalidQueryParamsMsg)
		return
	}

	changePwdErr := uh.Service.ChangeUserPwd(c, newPwd, username)
	if changePwdErr != nil {
		web.NewError(c, http.StatusInternalServerError, changePwdErr.Error())
		return
	}

	c.JSON(http.StatusOK, userResponse(http.StatusOK, config.UserPwdChangeMsg, nil))
}

func userResponse(status int, message string, data interface{}) *entity.UserResponse {
	return &entity.UserResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
