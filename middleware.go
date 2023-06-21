package main

import (
	"fmt"
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
		ulevel := session.Get("level").(uint8)
		if (ulevel >> 4) < level {
			fmt.Println(ulevel >> 4)
			ctx.AbortWithStatus(http.StatusForbidden) //403
			return
		}
	}
}
