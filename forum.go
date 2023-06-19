package main

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BrowseForumPost(ctx *gin.Context) {
	vpg := ctx.Param("page")
	pg, err := strconv.Atoi(vpg)
	if err != nil || pg < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	var forums []Forum
	var count int64 //总数,Count比rowsaffected更快(懒得用变量缓存了
	pg -= 1
	count, err = db.Count(&Forum{})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) //500,正常情况下不会出现
		log.Println(err)
		return
	}
	db.Limit(20, pg*20).Find(&forums)
	ctx.JSON(http.StatusOK, gin.H{
		"Body":      forums,
		"PageCount": math.Ceil(float64(count) / 20), //总页数
	})
	return
}

func BrowseUnitforumPost(ctx *gin.Context) {
	vmid := ctx.Param("mid")
	mid, err := strconv.Atoi(vmid)
	vpg := ctx.Param("page")
	pg, err := strconv.Atoi(vpg)
	var mainforum Forum
	has, err := db.ID(mid).Get(&mainforum)
	if err != nil || pg < 1 || has == false {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	var count int64
	pg -= 1
	var unitforum ForumComment
	count, err = db.Where("mid = ?", mid).Count(&unitforum)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) //500,正常情况下不会出现
		log.Println(err)
		return
	}
	var unitforums []ForumComment
	db.Limit(20, pg*20).Find(&unitforums)
	ctx.JSON(http.StatusOK, gin.H{
		"Body":      unitforums,
		"PageCount": math.Ceil(float64(count) / 20), //总页数
	})
}

func DBaddMainforum(title string, text string, author uint) {
	mainforum := &Forum{Title: title, Author: author, Views: 0}
	db.Insert(mainforum)
	unitforum := ForumComment{Text: text, Mid: uint(mainforum.Id), Author: author}
	db.Insert(unitforum)
	return
}

func DBaddUnitforum(text string, mid uint, author uint) {
	unitforum := ForumComment{Text: text, Mid: mid, Author: author}
	db.Insert(unitforum)
	return
}

func DBaddEmoji(emoji uint, uid uint) {
	var unitforum ForumComment
	db.ID(uid).Get(&unitforum)
	switch emoji {
	case 0:
		unitforum.Like++
	case 1:
		unitforum.Dislike++
	case 2:
		unitforum.Smile++
	case 3:
		unitforum.Firework++
	case 4:
		unitforum.Unhappy++
	case 5:
		unitforum.Heart++
	case 6:
		unitforum.Rocket++
	case 7:
		unitforum.Eyes++
	}
	db.ID(uid).Update(&unitforum)
	return
}
