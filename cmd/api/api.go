package main

import (
	"context"
	"net/http"
	"time"

	"github.com/KORLA2/SocialMedia/internal/store"
	"github.com/gin-gonic/gin"
)

type application struct {
	config config
	store  *store.Storage
}

type config struct {
	addr string
	db   dbConfig
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
	group.GET("health", app.HealthCheck)
	group.POST("posts", app.CreatePostHandler)
	group.POST("comments", app.CreateCommentHandler)

	middlewareUserGroup := group.Group("/user/:userID")
	middlewareUserGroup.Use(app.UsersContextMiddleWare)
	middlewareUserGroup.POST("/", app.CreateUserHandler)
	middlewareUserGroup.GET("/", app.GetUserHandler)
	middlewareUserGroup.PUT("follow", app.FollowUserHandler)
	middlewareUserGroup.PUT("unfollow", app.UnfollowUserHandler)
	middlewarePostGroup := group.Group("/posts/:postID")
	middlewarePostGroup.Use(app.PostsContextMiddleware)
	middlewarePostGroup.GET("/", app.GetPostHandler)
	middlewarePostGroup.DELETE("/", app.DeletePostHandler)
	middlewarePostGroup.PATCH("/", app.UpdatePostHandler)

	return router
}

func (app *application) Run(mux http.Handler) error {

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
