package main

import (
	"net/http"
	"strconv"

	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
)

type CommentPayload struct {
	Content string `json:"content" validate:"required,max=100"`
}

// PostComment godoc
//
//	@Summary		Creates a new comment on a post
//	@Description	Creates a valid comment on a post with content
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int				true	"Create comment for post ID"
//	@Param			payload	body		CommentPayload	true	"The comment content and metadata to create a new comment"
//	@Success		200		{object}	models.Comment
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/post/{postID}/comments [post]
func (a *application) CreateCommentHandler(c *gin.Context) {

	postID, _ := strconv.Atoi(c.Param("postID"))
	User, err := GetUserFromContext(c)
	if err != nil {
		a.InternalServerError(c, "User Context Not Found", err)
		return
	}

	ctx := c.Request.Context()
	var payload CommentPayload
	if err := c.BindJSON(&payload); err != nil {
		a.BadRequest(c, "Cannot Bind Json comments", err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		a.BadRequest(c, "Validation Failed", err)
		return
	}
	comment := models.Comment{

		PostID:  postID,
		UserID:  User.ID,
		Content: payload.Content,
	}

	if err := a.store.Comments.Create(ctx, &comment); err != nil {
		a.InternalServerError(c, "Could not create Comment", err)
		return
	}

	a.Success(c, "Comment Added Successfully", comment, http.StatusOK)
}
