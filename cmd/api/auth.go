package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/KORLA2/SocialMedia/internal/mailer"
	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func HashPassword(password string) string {

	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(pass)
}

func (a *application) RegisterUserHandler(c *gin.Context) {
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

	token := uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])
	log.Println(hashToken)

	User := models.User{

		Email:    payload.Email,
		Username: payload.Username,
		Password: payload.Password,
		Token:    token,
	}
	if err := a.store.Users.CreateAndInvite(ctx, &User, hashToken, a.config.mail.expiry); err != nil {
		a.BadRequest(c, "Cannot Create User", err)
		return
	}
	vars := struct {
		Username      string
		ActivationURL string
	}{
		User.Username,
		fmt.Sprintf("%s/confirm/%s", a.config.Frontend_URL, token),
	}

	if err := a.mailer.Send(mailer.UserWelcomeTemplateFile, User.Username, User.Email, vars, true); err != nil {

		// Rollback User Creation and Invitaion SAGA pattern
		if err := a.store.Users.Delete(ctx, User.ID); err != nil {
			a.InternalServerError(c, "User & Invite Transaction Deletion Failed", err)
		}

	}

	a.Success(c, "Created User Successfully", User, http.StatusOK)

}
func (a *application) ActivateUserHandler(c *gin.Context) {

	ctx := c.Request.Context()

	token := c.Param("token")
	log.Print(token)
	if err := a.store.Users.Activate(ctx, token); err != nil {
		a.BadRequest(c, "token error", err)
	}

}
