package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KORLA2/SocialMedia/internal/mailer"
	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginPayload struct {
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

func (a *application) CreateTokenHandler(c *gin.Context) {

	ctx := c.Request.Context()

	var payload LoginPayload

	if err := c.BindJSON(&payload); err != nil {

		a.BadRequest(c, "Cannot Bind Login Paylaod", err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		a.BadRequest(c, "Login Validation Failed", err)
		return
	}

	user := &models.User{}
	if err := a.store.Users.GetUserByUserName(ctx, payload.Username, user); err != nil {
		a.Unauthorized(c, "Username not found", err)
		return
	}

	// Token Creation

	claims := jwt.MapClaims{
		"sub": user.ID,
		"iss": a.config.auth.jwt.issuer,
		"aud": a.config.auth.jwt.audience,
		"exp": time.Now().Add(a.config.auth.jwt.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
	}
	token, err := a.auth.GenerateToken(claims)
	if err != nil {
		a.InternalServerError(c, "Cannot generate token", err)
		return
	}

	a.Success(c, "Successfully fetched User", token, http.StatusOK)

}

func (a *application) AuthenticateUserMiddleware(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	authSlice := strings.Split(authHeader, " ")

	if authSlice[0] != "Bearer" || len(authSlice) != 2 {
		a.Unauthorized(c, "Token Header Malformed", errors.New("need token to make any request"))
		return
	}

	token := authSlice[1]

	jwtToken, err := a.auth.ValidateToken(token)

	if err != nil {
		a.Unauthorized(c, "token is invalid", err)
		return
	}
	claims, _ := jwtToken.Claims.(jwt.MapClaims)

	userID, _ := strconv.Atoi(fmt.Sprintf("%v", claims["sub"]))
	ctx := c.Request.Context()
	user, err := a.store.Users.GetUserByID(ctx, userID)
	if err != nil {
		a.Unauthorized(c, "token is invalid", err)
		return
	}

	c.Set("User", user)

}
