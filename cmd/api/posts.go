package main

import (
	"context"
	"errors"
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

// CreatePost godoc
//
//	@Summary		Creates a new post
//	@Description	Creates a valid user post with title, content, and tags
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		PostPayload	true	"Post data"
//	@Success		200		{object}	models.Post
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (a *application) CreatePostHandler(c *gin.Context) {

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
	User, _ := GetUserFromContext(c)

	post := models.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  User.ID,
		Tags:    payload.Tags,
	}

	if err := a.store.Posts.Create(ctx, &post); err != nil {
		a.InternalServerError(c, "Cannot Create Post", err)

	}
	a.Success(c, "Post Created Successfully", post, http.StatusCreated)

}

// GetPost           godoc
//
//	@Summary		Fetches a user post
//	@Description	Fetches any valid and verified  user's post  by ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"postID"
//	@Success		200		{object}	models.Post
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [get]
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

// Delete Post     godoc
//
//	@Summary		Deletes a user post
//	@Description	Deletes a user's post by ID Admin can delete any post and user can delete his own post.
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"postID"
//	@Success		200		{object}	models.Post
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [delete]
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

// UpdatePost           godoc
//
//	@Summary		Updates a user post
//	@Description	Updates user's post  by ID Admin and moderator can update any user post, user can update his own post
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"postID"
//	@Param			payload	body		UpdatePostPayload	true	"Post data"
//	@Success		200		{object}	models.Post
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [patch]
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

	if !exists {
		return nil, fmt.Errorf("post context not fetched")
	}
	post := postInterface.(*models.Post)

	return post, nil

}

func (a *application) CheckPostOwnership(minimumRequiredRole string, handler gin.HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		User, _ := GetUserFromContext(c)
		Post, _ := getPostFromCtx(c)

		if Post.UserID == User.ID {
			handler(c)
			return
		}

		allowed, _ := a.checkPrecedence(ctx, minimumRequiredRole, User.Role.Level)

		if !allowed {
			a.ForbiddenError(c, "Not Allowed", errors.New("you are not allowed to perform this operation"))
			return
		}
		handler(c)

	}
}

func (a *application) checkPrecedence(ctx context.Context, minimumRequiredRole string, myLevel int) (bool, error) {

	role, err := a.store.Roles.GetRoleByName(ctx, minimumRequiredRole)
	if err != nil {
		return false, err
	}
	return myLevel >= role.Level, nil
}
