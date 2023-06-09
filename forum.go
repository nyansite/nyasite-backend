package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"log"
	"math"
)

func BrowseForumPost(ctx *gin.Context)  {
	vpg := ctx.Param("page")
	pg, err := strconv.Atoi(vpg)
	if err != nil || pg < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	var forums []Video
	var count int64 //总数,Count比rowsaffected更快(懒得用变量缓存了
	pg -= 1
	count, err = db.Count(&Forum{})
	if err != nil{
		ctx.AbortWithStatus(http.StatusInternalServerError) //500,正常情况下不会出现
		log.Println(err)
		return
	}
	db.Limit(20, pg*20).Find(&forums)
	ctx.JSON(http.StatusOK, gin.H{
		"Body":      forums,
		"PageCount": math.Ceil(float64(count) / 20), //总页数
	})
}