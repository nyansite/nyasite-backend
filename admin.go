package main

import (
	"fmt"
	"math"
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	UUID "github.com/google/uuid"
)

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
	if err != nil{
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
func UploadVideo(c *gin.Context) {
	title := c.PostForm("Title")
	description := c.PostForm("Description") //简介
	f, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}

	uuid := UUID.New()
	fpath := "./temporary/" + uuid.String() + path.Ext(f.Filename)
	c.SaveUploadedFile(f, fpath)
	go SaveVideo(fpath, title, description, uuid.String())
}
