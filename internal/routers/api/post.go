package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"paopao-ce-teaching/internal/services"
	"paopao-ce-teaching/pkg/app"
	"paopao-ce-teaching/pkg/errors"
)

func CreatePost(c *gin.Context) {
	param := services.PostCreationReq{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		log.Errorf("app.BindAndValid errs: %v", errs)
		app.ToErrorResponse(c, errors.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	userID, _ := c.Get("UID")
	post, err := services.CreatePost(c, userID.(int64), param)

	if err != nil {
		log.Errorf("service.CreatePost err: %v\n", err)
		app.ToErrorResponse(c, errors.CreatePostFailed)
		return
	}

	app.ToResponse(c, post)
}
