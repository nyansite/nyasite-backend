package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func PrivilegeLevel(level uint8) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		is_login, err := ctx.Cookie("is_login")
		if is_login != "true" || err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized) //未登录返回401
			return
		}

		userid := session.Get("userid")
		var user User
		if has, _ := db.ID(userid).Get(&user); has == false { //用户不存在
			ctx.Status(http.StatusBadRequest)
			return
		}
		if session.Get("pwd-8") != string(user.Passwd[:8]) {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ulevel := session.Get("level").(uint8)
		if (ulevel >> 4) < level {
			ctx.AbortWithStatus(http.StatusForbidden) //403
			return
		}
	}
}

// github.com/gin-gonic/gin/blob/44d0dd70924dd154e3b98bc340accc53484efa9c/logger.go#L134
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Request.Host+param.Path, //加上了host
		param.ErrorMessage,
	)
}
