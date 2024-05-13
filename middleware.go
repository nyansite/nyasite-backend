package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
		if has, _ := db.ID(int(userid.(int64))).Get(&user); !has { //用户不存在,可能是由于账户被删除
			ctx.AbortWithStatus(http.StatusUnauthorized) //401
			return
		}
		if string(session.Get("pwd-8").([]byte)) != string(user.Passwd[:8]) { //确保改密码会导致所有登陆设备登陆失败
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		ulevel := user.Level
		if (ulevel >> 4) < level {
			ctx.AbortWithStatus(http.StatusForbidden) //403
			return
		}
	}
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