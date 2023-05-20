package main

import (
	_ "fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func NewTag(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("is_login") != true {
		c.AbortWithStatus(http.StatusUnauthorized) //返回401
		return
	}
	level := session.Get("level").(uint)
	privilege_level := level >> 4
	if privilege_level < 10 {
		c.AbortWithStatus(http.StatusForbidden) //403
		return
	}
	tagname := c.PostForm("tagname")
	if db.First(&TagText{}, "Text = ?", tagname).RowsAffected != 0 {
		c.AbortWithStatus(StatusRepeatTag)
		return
	}
	db.Create(&TagText{Text: tagname})
	c.AbortWithStatus(http.StatusOK)
}
