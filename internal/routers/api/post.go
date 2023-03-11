package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"paopao-ce-teaching/internal/services"
	"paopao-ce-teaching/pkg/app"
	"paopao-ce-teaching/pkg/convert"
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

func GetPost(c *gin.Context) {
	postID := convert.StrTo(c.Query("id")).MustInt64()
	postFormatted, err := services.GetPost(postID)

	if err != nil {
		log.Errorf("service.GetPost err: %v\n", err)
		app.ToErrorResponse(c, errors.GetPostFailed)
		return
	}

	app.ToResponse(c, postFormatted)
}

func DeletePost(c *gin.Context) {
	param := services.PostDelReq{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		log.Errorf("app.BindAndValid errs: %v", errs)
		app.ToErrorResponse(c, errors.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	userID, exist := c.Get("UID")
	if !exist {
		app.ToErrorResponse(c, errors.NoPermission)
		return
	}

	err := services.DeletePost(userID.(int64), param.ID)
	if err != nil {
		log.Errorf("service.DeletePost err: %v\n", err)
		app.ToErrorResponse(c, err)
		return
	}

	app.ToResponse(c, nil)
}
