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
	Level    int    `json:"level" validate:"required,min=1,max=3"`
}

type LoginUserPayload struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=4,max=100"`
}

func HashPassword(password string) string {

	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(pass)
}

func ValidateUserPassword(payloadpassword, userpassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(userpassword), []byte(payloadpassword))
	if err != nil {
		return false
	}

	return true

}

// User SignUp           godoc
//
//	@Summary		Registers a new user
//	@Description	Registers a new user with email, username, password, and role level: 1 for user 2 for moderator 3 for admin
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UserPayload	true	"User Signup"
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error

// @Router	/authenticate/user/signup [post]
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
		Role:     models.Role{Level: payload.Level},
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

// UserActivation godoc
//
//	@Summary		Activates a new user
//	@Description	Activates a new user account via the token sent to their email
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string				true	"Activation token"
//
//	@Success		200		{object}	map[string]string	"User activated successfully"
//
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error

// @Router	/authenticate/user/activate/{token} [put]
func (a *application) ActivateUserHandler(c *gin.Context) {

	ctx := c.Request.Context()

	token := c.Param("token")
	log.Print(token)
	if err := a.store.Users.Activate(ctx, token); err != nil {
		a.BadRequest(c, "token error", err)
		return
	}
	a.Success(c, "User activated Successfully", "User activated Successfully", http.StatusOK)
}

// User Login godoc
//
//	@Summary		Logs in a user and returns a JWT token
//	@Description	Logs in a user with username and password, returning a JWT token for authenticated requests
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		LoginUserPayload	true	"User Login Payload"
//
//	@Success		200		{object}	map[string]string	"User activated successfully"
//
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error

// @Router	/authenticate/user/signin [post]
func (a *application) LoginUserHandler(c *gin.Context) {

	ctx := c.Request.Context()

	var payload LoginUserPayload

	if err := c.BindJSON(&payload); err != nil {

		a.BadRequest(c, "Cannot Bind Login Paylaod", err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		a.BadRequest(c, "Login Validation Failed", err)
		return
	}
	user, err := a.store.Users.GetUserByUserName(ctx, payload.Username)
	if err != nil {
		a.Unauthorized(c, "Username or password is incorrect", err)
		return
	}

	if ok := ValidateUserPassword(payload.Password, user.Password); !ok {
		a.Unauthorized(c, "Username or password is incorrect", err)
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
