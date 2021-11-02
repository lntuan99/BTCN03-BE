package api_classroom

import (
	"web2/btcn/api/base"
	"web2/btcn/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerGetClassroomList(c *gin.Context) {
	success, status, data := methods.MethodGetClassroomList(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerCreateClassroom(c *gin.Context) {
	success, status, data := methods.MethodCreateClassroom(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

