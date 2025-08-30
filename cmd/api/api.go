package main

import (
	"context"
	"net/http"
	"time"

	"github.com/KORLA2/SocialMedia/cmd/docs"
	"github.com/KORLA2/SocialMedia/internal/auth"
	"github.com/KORLA2/SocialMedia/internal/mailer"
	"github.com/KORLA2/SocialMedia/internal/store"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type application struct {
	config config
	store  *store.Storage
	mailer mailer.Client
	auth   auth.Authenticator
}

type config struct {
	addr         string
	db           dbConfig
	Frontend_URL string
	mail         mailConfig
	auth         authConfig
}

type authConfig struct {
	basic basicConfig
	jwt   jwtConfig
}
type basicConfig struct {
	user     string
	password string
}

type jwtConfig struct {
	secret   string
	audience string
	issuer   string
	exp      time.Duration
}

type mailConfig struct {
	sendgrid  SendGridConfig
	FromEmail string
	expiry    time.Duration
}

type SendGridConfig struct {
	API_KEY string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(RequestTimeOut(30 * time.Second))
	group := router.Group("/api/v1")
	group.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1)))

	group.GET("health", gin.BasicAuth(gin.Accounts{
		app.config.auth.basic.user: app.config.auth.basic.password,
	}), app.HealthCheck)

	group.POST("posts", app.CreatePostHandler)
	group.POST("comments", app.CreateCommentHandler)
	group.GET("/user/feed", app.GetUserFeedHandler)
	middlewareUserGroup := group.Group("/user/:userID")
	middlewareUserGroup.Use(app.AuthenticateUserMiddleware)
	middlewareUserGroup.Use(app.UsersContextMiddleWare)
	middlewareUserGroup.GET("/", app.GetUserHandler)
	middlewareUserGroup.PUT("follow", app.FollowUserHandler)
	middlewareUserGroup.PUT("unfollow", app.UnfollowUserHandler)

	middlewarePostGroup := group.Group("/posts/:postID")
	middlewarePostGroup.Use(app.AuthenticateUserMiddleware)
	middlewarePostGroup.Use(app.PostsContextMiddleware)
	middlewarePostGroup.GET("/", app.GetPostHandler)
	middlewarePostGroup.DELETE("/", app.DeletePostHandler)
	middlewarePostGroup.PATCH("/", app.UpdatePostHandler)

	// Public Routes
	middlewareAuthGroup := group.Group("/authenticate/user")

	middlewareAuthGroup.POST("/", app.RegisterUserHandler)
	middlewareAuthGroup.POST("token", app.CreateTokenHandler)
	middlewareAuthGroup.PUT("activate/:token", app.ActivateUserHandler)

	return router
}

func (app *application) Run(mux http.Handler) error {
	docs.SwaggerInfo.BasePath = "/api/v1"
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	return server.ListenAndServe()

}

func RequestTimeOut(time time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), time)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan bool)
		go func() {

			c.Next()
			close(done)
		}()
		select {
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "Request Timed Out",
			})

		case <-done:

		}

	}

}
