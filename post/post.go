package post

import (
	"cute_site/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"encoding/json"
	"fmt"
	"strings"
)

func CreatNewPost(db *gorm.DB, c *gin.Context) { //发帖
	var newMainPost models.MainPost
	var user models.User
	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	user_p := body["user_id"]
	content := body["content"]
	title := body["title"]
	result := db.Last(&user, "id = ?", user_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(612)
		return
	}
	db.AutoMigrate(&models.MainPost{})
	db.Create(&models.MainPost{User_p: user_p, Title: title, Video_p: "indepent", Views: 0, Likes: 0, ContentShow: content})
	db.Last(&newMainPost)
	mainpost_p := fmt.Sprintf("%d", newMainPost.ID)
	SendMain(db, c, mainpost_p, user_p, content)
}

func Send(db *gorm.DB, c *gin.Context) { //跟帖
	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	mainPost_p := body["mainpost_id"]
	user_p := body["user_id"]
	content := body["content"]
	SendMain(db, c, mainPost_p, user_p, content)
}

func SendMain(db *gorm.DB, c *gin.Context, mainPost_p string, user_p string, content string) { //无论是发新帖还是根贴都要用到这个函数
	var lastPost models.UnitPost
	var user models.User
	var mainPost models.MainPost
	result := db.Last(&user, "id = ?", user_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(612)
		return
	}
	db.AutoMigrate(&models.MainPost{})
	result = db.Last(&mainPost, "id = ?", mainPost_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(622)
		return
	}
	db.AutoMigrate(&models.UnitPost{})
	db.Last(&lastPost, "mainpost_p = ?", mainPost_p)
	db.Create(&models.UnitPost{MainPost_p: mainPost_p, User_p: user_p, Content: content})
	c.AbortWithStatus(621)
	return
}

func CommentSend(db *gorm.DB, c *gin.Context) {
	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	post_p := body["post_p"]
	content := body["content"]
	user_p := body["user_p"]
	var user models.User
	var unitPost models.UnitPost
	result := db.Last(&user, "id = ?", user_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(612)
		return
	}
	result = db.Last(&unitPost, "id = ?", post_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(622)
		return
	}
	db.Create(&models.Comment{Content: content, Post_p: post_p, User_p: user_p})
	c.AbortWithStatus(621)
	return
}

func GetMainPost(db *gorm.DB, c *gin.Context) {
	mainPost_p := c.Request.Header.Get("mainpost_p")
	var mainpost models.MainPost
	var unitposts []models.UnitPost
	result := db.Last(&mainpost, "id = ?", mainPost_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(622)
		return
	}
	db.Last(&unitposts, "mainpost_p = ?", mainPost_p)
	unitPostIDList := ""
	for _, value := range unitposts {
		strings.Join([]string{unitPostIDList, fmt.Sprintf("%d", value.ID), ","}, "")
	}
	c.JSON(http.StatusOK, gin.H{
		"title":          mainpost.Title,
		"views":          mainpost.Views,
		"user_p":         mainpost.User_p,
		"utilpostidlist": unitPostIDList,
	})
	mainpost.Views = mainpost.Views + 1
	return
}

func GetPost(db *gorm.DB, c *gin.Context) {
	post_p := c.Request.Header.Get("post_p")
	var unitpost models.UnitPost
	var comments []models.Comment
	result := db.Last(&unitpost, "id = ?", post_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(622)
		return
	}
	db.Find(&comments, "post_p = ?", post_p)
	commentIdList := ""
	for _, value := range comments {
		strings.Join([]string{commentIdList, fmt.Sprintf("%d", value.ID), ","}, "")
	}
	c.JSON(http.StatusOK, gin.H{
		"content":       unitpost.Content,
		"user_p":        unitpost.User_p,
		"commentidlist": commentIdList,
	})
	return
}

func GetComment(db *gorm.DB, c *gin.Context) {
	comment_p := c.Request.Header.Get("comment_p")
	var comment models.Comment
	result := db.Last(&comment, "id = ?", comment_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(622)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"content": comment.Content,
		"user_p":  comment.User_p,
	})
	return
}
