package user

import (
	"go-manage-hex/cmd/config"
	"go-manage-hex/internal/app/user"
	"net/http"
	"strconv"

	entity "go-manage-hex/internal/core/user"
	dto "go-manage-hex/internal/infrastructure/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/gustyaguero21/go-core/pkg/web"
)

type UserHandler struct {
	Service     user.Usecases
	AuthService entity.Authorization
}

func NewUserHandler(service user.Usecases, auth entity.Authorization) *UserHandler {
	return &UserHandler{
		Service:     service,
		AuthService: auth,
	}
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

	var dto dto.UpdateDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		web.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	user := entity.User{
		Name:     dto.Name,
		LastName: dto.LastName,
		Email:    dto.Email,
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

func (uh *UserHandler) LoginUser(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var dto dto.LoginRequestDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		web.NewError(c, http.StatusBadRequest, config.InvalidBodyMsg)
		return
	}

	user := entity.User{
		Username: dto.Username,
		Password: dto.Password,
	}

	if loginErr := uh.Service.Login(c, user.Username, user.Password); loginErr != nil {
		web.NewError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := uh.AuthService.GenerateJWT(user.Username, user.Password)
	if err != nil {
		web.NewError(c, http.StatusInternalServerError, "error generating token")
		return
	}

	c.JSON(http.StatusOK, userResponse(http.StatusOK, "user logged", token))
}

func userResponse(status int, message string, data interface{}) *dto.UserResponseDTO {
	return &dto.UserResponseDTO{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
