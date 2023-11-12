package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func PrivilegeLevel(level uint8) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if session.Get("is_login") != true {
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
