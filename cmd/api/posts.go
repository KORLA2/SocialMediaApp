package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type PostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}
type UpdatePostPayload struct {
	Title   string `json:"title" validate:"omitempty,max=100"`
	Content string `json:"content" validate:"omitempty,max=1000"`
}

func (a *application) CreatePostHandler(c *gin.Context) {

	// ctx
	ctx := c.Request.Context()
	var payload PostPayload

	if err := c.BindJSON(&payload); err != nil {
		a.BadRequest(c, "cannotBind Json", err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		a.BadRequest(c, "Validation Failed", err)
		return
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

	post, err := getPostFromCtx(c)

	if err != nil {
		a.InternalServerError(c, "Unable to fetch the post COntext", err)
		return
	}
	comments, err := a.store.Comments.GetCommentsByPostID(c.Request.Context(), (*post).ID)
	if err != nil {
		a.InternalServerError(c, "Couldn't get comments for this post", err)

		return
	}
	(*post).Comments = comments

	a.Success(c, "Post Fetched Successfully", *post, http.StatusOK)

}

func (a *application) DeletePostHandler(c *gin.Context) {
	postIDstring := c.Param("postID")
	postID, _ := strconv.Atoi(postIDstring)
	ctx := c.Request.Context()

	if err := a.store.Posts.DeletePostByID(ctx, postID); err != nil {
		a.InternalServerError(c, "Could Not Delete Post", err)
		return
	}

	a.Success(c, "Post Deleted Successfully", postID, http.StatusOK)

}

func (a *application) UpdatePostHandler(c *gin.Context) {

	post, err := getPostFromCtx(c)

	if err != nil {
		a.InternalServerError(c, "Unable to fetch the post COntext", err)
		return
	}

	var payload UpdatePostPayload

	if err := c.BindJSON(&payload); err != nil {
		a.BadRequest(c, "Cannot Bind Json while updating post", err)
		return
	}

	if payload.Content != "" {

		post.Content = payload.Content
	}
	if payload.Title != "" {
		post.Title = payload.Title
	}
	log.Println(payload)

	if err := a.store.Posts.UpdatePostByID(c.Request.Context(), post); err != nil {
		a.InternalServerError(c, "Cannot Update Post", err)
		return
	}

	a.Success(c, "Successfully Updated Post", *post, http.StatusOK)
}

func (a *application) PostsContextMiddleware(c *gin.Context) {

	postIDstring := c.Param("postID")
	postID, _ := strconv.Atoi(postIDstring)

	ctx := c.Request.Context()
	post, err := a.store.Posts.GetPostByID(ctx, postID)
	if err != nil {

		a.InternalServerError(c, "Couldn't get post", err)

		return
	}
	c.Set("post", post)

}

func getPostFromCtx(c *gin.Context) (*models.Post, error) {

	postInterface, exists := c.Get("post")

	if exists != true {
		return nil, fmt.Errorf("post context not fetched")
	}
	post := postInterface.(*models.Post)

	return post, nil

}
