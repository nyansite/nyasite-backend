package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AdminCheck() gin.HandlerFunc {
	return func(c *gin.Context){
		session := sessions.Default(c)
		if session.Get("is_login") != true {
			c.AbortWithStatus(http.StatusUnauthorized) //返回401
			return
		}
		level := session.Get("level").(uint8)
		privilege_level := level >> 4
		if privilege_level < 15 { //15满级权限才能进入
			c.String(http.StatusForbidden, "梦里啥都有")	//403
			c.Abort()
			return
		}
	}
}