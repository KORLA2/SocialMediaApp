package main

import (
	"net/http"
	"strconv"

	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
var validate = validator.New()
type PostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required, max=1000"`
	Tags    []string `json:"tags"`
}

func (a *application) CreatePostHandler(c *gin.Context) {

	// ctx
	ctx := c.Request.Context()
	var payload PostPayload


	if err := c.BindJSON(&payload); err != nil {
		a.BadRequest(c, "cannotBind Json", err)
		return;
	}

	if err:= validate.Struct(payload); err != nil {
		a.BadRequest(c, "Validation Failed", err)
			return;
	}
	post := models.Post{
		Title:   payload.Title,
		Content: payload.Content,
		User_ID: 1,
		Tags:    payload.Tags,
	}

	if err := a.store.Posts.Create(ctx, &post); err != nil {
		a.InternalServerError(c, "Cannot Create Post", err)

	}
	a.Success(c, "Post Created Successfully", post, http.StatusCreated)

}

func (a *application) GetPostHandler(c *gin.Context) {

	postIDstring := c.Param("postID")
	postID, _ := strconv.Atoi(postIDstring)

	ctx := c.Request.Context()
	post, err := a.store.Posts.GetPostByID(ctx, postID)
	if err != nil {

		a.InternalServerError(c, "Couldn't get post", err)

		return
	}
	a.Success(c, "Post Fetched Successfully", *post, http.StatusOK)

}
