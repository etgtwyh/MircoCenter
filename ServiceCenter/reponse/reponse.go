package reponse

import (
	"fmt"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/execption"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// 统一HTTP返回消息格式
func RepFailed(ctx *gin.Context, err error) {
	//	断言err判断是否是Exception
	if target, ok := err.(*execption.Exception); ok {
		ctx.JSON(target.Code, target)
		ctx.Abort()
	} else if _, ok := err.(validator.ValidationErrors); ok {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Request parameter error:", err.Error()))
		ctx.Abort()
	} else {
		ctx.JSON(http.StatusInternalServerError, "unknown error")
		ctx.Abort()
	}
}

func RepSuccess(ctx *gin.Context, format any) {
	ctx.JSON(http.StatusOK, format)
}
