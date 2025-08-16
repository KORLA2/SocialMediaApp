package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"user_name" validate:"required"`
	Password string `json:"-"`
}

func HashPassword(password string) string {

	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(pass)
}

func (a *application) CreateUserHandler(c *gin.Context) {

	ctx := c.Request.Context()

	var payload UserPayload
	if err := c.BindJSON(&payload); err != nil {

		a.BadRequest(c, "Cannot Bind User Json", err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		a.BadRequest(c, "User Validation Failed", err)
		return
	}

	payload.Password = HashPassword(payload.Password)

	User := models.User{

		Email:    payload.Email,
		Username: payload.Username,
		Password: payload.Password,
	}
	if err := a.store.Users.Create(ctx, &User); err != nil {
		a.InternalServerError(c, "Cannot Create User", err)
		return
	}

	a.Success(c, "Created User Successfully", User, http.StatusOK)

}

func (a *application) GetUserHandler(c *gin.Context) {

	user, err := GetUserFromContext(c)
	if err != nil {
		a.InternalServerError(c, "User Context Not set", err)
		return
	}

	a.Success(c, "Fecthed User Successfully ", *user, http.StatusOK)

}

func (a *application) UsersContextMiddleWare(c *gin.Context) {

	ctx := c.Request.Context()

	userIDstring := c.Param("userID")
	userID, _ := strconv.Atoi(userIDstring)
	user, err := a.store.Users.GetUserByID(ctx, userID)

	if err != nil {
		a.InternalServerError(c, "Cannot Get User", err)
		return
	}
	c.Set("User", user)

}

func GetUserFromContext(ctx *gin.Context) (*models.User, error) {

	usersInterface, exists := ctx.Get("User")

	if !exists {

		return nil, fmt.Errorf("user context not set")

	}

	return usersInterface.(*models.User), nil

}
