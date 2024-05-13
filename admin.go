package main

import (
	"fmt"

	"math"
	"net/http"
	"strconv"

	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func NewTag(c *gin.Context) {
	is_login, _ := c.Cookie("is_login")
	if is_login != "true" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tagname := c.PostForm("tagname")
	if has, _ := db.Exist(&TagModel{Text: tagname}); has{
		c.AbortWithStatus(http.StatusConflict)//409 请求与当前内容冲突
		return
	}
	db.Insert(&TagModel{Text: tagname})
	c.AbortWithStatus(http.StatusOK)
}

func AdminVideoPost(ctx *gin.Context) {
	vpg := ctx.Param("page")
	pg, err := strconv.Atoi(vpg)
	if err != nil || pg < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	var videos []Video
	var count int64 //总数,Count比rowsaffected更快(懒得用变量缓存了
	pg -= 1
	count, err = db.Count(&Video{})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) //500,正常情况下不会出现
		fmt.Println(err)
		return
	}
	db.Limit(20, pg*20).Find(&videos)
	ctx.JSON(http.StatusOK, gin.H{
		"Body":      videos,
		"PageCount": math.Ceil(float64(count) / 20), //总页数
	})
}

// 不审核直接上传,测试接口

func GetSessionSecret() [][]byte {

	return nil
}
