package main

import(
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"net/http"
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
		if privilege_level <= 15 { //15满级权限才能进入
			c.AbortWithStatus(http.StatusForbidden) //403
			return
		}
	}
}