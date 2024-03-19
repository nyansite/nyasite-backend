package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CheckPrivilege(level uint8) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		is_login, _ := ctx.Cookie("is_login")
		if is_login != "true" {
			ctx.AbortWithStatus(http.StatusUnauthorized) //未登录返回401
			return
		}
		userid := session.Get("userid")
		var user User
		if has, _ := db.ID(userid).Get(&user); !has { //用户不存在,可能是由于账户被删除
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if session.Get("pwd-8") != string(user.Passwd[:8]) {//确保改密码会导致所有登陆设备登陆失败
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		ulevel := user.Level
		if (ulevel >> 4) < level {
			ctx.AbortWithStatus(http.StatusForbidden) //403
			return
		}
		ctx.Next()
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

// 限制请求大小(不是文件大小),能够节约cpu,似乎并不能节约网络
func LimitRequestBody(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize) //创建reader并不会立即读取,读取是 reader.Read() 或其他
		err := c.Request.ParseMultipartForm(maxSize)                            //尝试读取Request.Body
		if err != nil {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
			return
		}
		c.Next()
	}
}
