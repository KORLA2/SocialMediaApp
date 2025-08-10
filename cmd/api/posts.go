package main

import (
	"net/http"
	"strconv"

	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
)

type PostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (a *application) CreatePostHandler(c *gin.Context) {

	// ctx
	ctx := c.Request.Context()
	var payload PostPayload

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Cannot Bind JSON": err.Error(),
		})
	}
	post := models.Post{
		Title:   payload.Title,
		Content: payload.Content,
		User_ID: 1,
		Tags:    payload.Tags,
	}

	if err := a.store.Posts.Create(ctx, &post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"Sucess Post Created": post,
	})

}

func (a *application) GetPostHandler(c *gin.Context) {

	postIDstring := c.Param("postID")
	postID, _ := strconv.Atoi(postIDstring)

	ctx := c.Request.Context()
	post, err := a.store.Posts.GetPostByID(ctx, postID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"Couldn't get post": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"PostFound": *post,
	})

}
