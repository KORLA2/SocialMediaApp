package main

import (
	"net/http"

	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
)

type CommentPayload struct {
	Content string `json:"content" validate:"required,max=100"`
}

func (a *application) CreateCommentHandler(c *gin.Context) {

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

		PostID:  7,
		UserID:  1,
		Content: payload.Content,
	}

	if err := a.store.Comments.Create(ctx, &comment); err != nil {
		a.InternalServerError(c, "Could not create Comment", err)
		return
	}

	a.Success(c, "Comment Added Successfully", comment, http.StatusOK)
}
