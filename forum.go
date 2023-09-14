package main

import (
	"github.com/gin-contrib/sessions"

	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"fmt"
)

func BrowseAllForumPost(ctx *gin.Context) {
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

func BrowseForumPost(ctx *gin.Context) {
	vkind := ctx.Param("board")
	kind, err := strconv.Atoi(vkind)
	vpg := ctx.Param("page")
	pg, err := strconv.Atoi(vpg)
	var chose []int
	switch kind {
	case 0:
		chose = append(chose, 0)
	case 1:
		chose = append(chose, 1, 2)
	case 2:
		chose = append(chose, 3, 4)
	case 3:
		chose = append(chose, 5)
	default:
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var forums []Forum
	var count int64 //总数,Count比rowsaffected更快(懒得用变量缓存了
	pg -= 1
	count, err = db.In("kind", chose).Count(&Forum{})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) //500,正常情况下不会出现
		log.Println(err)
		return
	}
	db.In("kind", chose).Limit(20, pg*20).Find(&forums)
	ctx.JSON(http.StatusOK, gin.H{
		"Body":      forums,
		"PageCount": math.Ceil(float64(count) / 20), //总页数
	})
	return

}

func BrowseUnitforumPost(ctx *gin.Context) {
	session := sessions.Default(ctx)
	author := session.Get("userid")
	vauthor, _ := author.(int)
	uauthor := uint(vauthor)
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
	count, err = db.In("mid", mid).Count(&unitforum)
	var unitforums []ForumComment
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) //500,正常情况下不会出现
		log.Println(err)
		return
	}
	db.In("mid", mid).Limit(20, pg*20).Find(&unitforums)
	var unitforumsReturn []ForumComment
	for _, i := range unitforums {
		var emojiRecord EmojiRecord
		count, _ := db.Where("author = ? AND uid = ?", author, i.Id).Count(&EmojiRecord{})
		if count == 0 {
			i.Choose = 0
		} else {
			db.Where("author = ? AND uid = ?", uauthor, i.Id).Get(&emojiRecord)
			i.Choose = emojiRecord.Emoji + 1
		}
		i.Like--
		i.Dislike--
		i.Smile--
		i.Celebration--
		i.Confused--
		i.Heart--
		i.Rocket--
		i.Eyes--
		unitforumsReturn = append(unitforumsReturn, i)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Origin":    mainforum,
		"Body":      unitforumsReturn,
		"PageCount": math.Ceil(float64(count) / 20), //总页数
	})
}

//因为xorm 值为0默认不update 所以表情从1开始计数

func DBaddMainforum(title string, text string, author int, kind int) {
	mainforum := &Forum{Title: title, Author: author, Views: 0, Kind: kind}
	db.Insert(mainforum)
	unitforum := ForumComment{Text: text, Mid: int(mainforum.Id), Author: author}
	db.Insert(unitforum)
	return
}

func DBaddUnitforum(text string, mid int, author int) {
	unitforum := ForumComment{Text: text, Mid: mid, Author: author, Like: 1, Dislike: 1, Smile: 1, Celebration: 1, Confused: 1, Heart: 1, Rocket: 1, Eyes: 1}
	db.Insert(unitforum)
	return
}

func DBaddEmoji(emoji int, uid int, author int) {
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
		unitforum.Celebration++
	case 4:
		unitforum.Confused++
	case 5:
		unitforum.Heart++
	case 6:
		unitforum.Rocket++
	case 7:
		unitforum.Eyes++
	}
	emojiRecord := EmojiRecord{Author: author, Uid: uid, Emoji: emoji}
	db.Insert(&emojiRecord)
	db.ID(uid).Update(&unitforum)
	return
}

func DBchangeEmoji(emoji int, uid int, author int) {
	var unitforum ForumComment
	var emojiRecord EmojiRecord
	db.ID(uid).Get(&unitforum)
	db.Where("author = ? and uid = ?", author, uid).Get(&emojiRecord)
	fmt.Println(unitforum)
	fmt.Println(emojiRecord)
	oEmoji := emojiRecord.Emoji
	switch oEmoji {
	case 0:
		unitforum.Like--
	case 1:
		unitforum.Dislike--
	case 2:
		unitforum.Smile--
	case 3:
		unitforum.Celebration--
	case 4:
		unitforum.Confused--
	case 5:
		unitforum.Heart--
	case 6:
		unitforum.Rocket--
	case 7:
		unitforum.Eyes--
	}
	switch emoji {
	case 0:
		unitforum.Like++
	case 1:
		unitforum.Dislike++
	case 2:
		unitforum.Smile++
	case 3:
		unitforum.Celebration++
	case 4:
		unitforum.Confused++
	case 5:
		unitforum.Heart++
	case 6:
		unitforum.Rocket++
	case 7:
		unitforum.Eyes++
	}
	emojiRecord.Emoji = emoji
	fmt.Println(unitforum)
	fmt.Println(emojiRecord)
	db.ID(uid).Update(&unitforum)
	db.Where("author = ? and uid = ?", author, uid).Cols("emoji").Update(&emojiRecord)
	return
}

func AddMainforum(ctx *gin.Context) {
	session := sessions.Default(ctx)
	title, text, kind := ctx.PostForm("title"), ctx.PostForm("text"), ctx.PostForm("type")
	vkind, _ := strconv.Atoi(kind)
	var ukind uint8
	switch vkind {
	case 0:
		ukind = 1
	case 1:
		ukind = 3
	case 2:
		ukind = 5
	default:
		ctx.AbortWithStatus(http.StatusBadRequest) //传入了错误的分区数据
		return
	}
	author := session.Get("userid")
	vauthor := author.(int64)
	uauthor := int(vauthor)
	DBaddMainforum(title, text, uauthor, int(ukind))
	return
}

func AddUnitforum(ctx *gin.Context) {
	session := sessions.Default(ctx)
	author := session.Get("userid")
	vauthor := author.(int64)
	uauthor := int(vauthor)
	mid, text := ctx.PostForm("mid"), ctx.PostForm("text")
	vmid, _ := strconv.Atoi(mid)
	umid := int(vmid)
	DBaddUnitforum(text, umid, uauthor)
	return
}

func AddEmoji(ctx *gin.Context) {
	session := sessions.Default(ctx)
	author := session.Get("userid")
	vauthor := author.(int64)
	uauthor := int(vauthor)
	emoji, uid := ctx.PostForm("emoji"), ctx.PostForm("uid")
	vuid, _ := strconv.Atoi(uid)
	vemoji, _ := strconv.Atoi(emoji)
	exist, _ := db.Where("author = ? and uid = ?", uauthor, uid).Count(&EmojiRecord{})
	if exist != 0 {
		DBchangeEmoji(vemoji, vuid, uauthor)
		return
	}
	if vemoji > 7 {
		ctx.AbortWithStatus(http.StatusBadRequest) //传入的表情编号>7(不存在)
		return
	}
	DBaddEmoji(vemoji, vuid, uauthor)
	return
}

func FinishForum(ctx *gin.Context) {
	session := sessions.Default(ctx)
	mid := ctx.PostForm("mid")
	vmid, _ := strconv.Atoi(mid)
	author := session.Get("userid")
	vauthor := author.(int64)
	uauthor := int(vauthor)
	var mainforum Forum
	db.ID(vmid).Get(&mainforum)
	if uauthor == mainforum.Author {
		mainforum.Kind++
	} else {
		ctx.AbortWithStatus(http.StatusBadRequest) //如果不是贴主不能完结
	}
	db.ID(vmid).Update(&mainforum)
	return
}
