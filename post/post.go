package post

import (
	"cute_site/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatNewPost(db *gorm.DB, c *gin.Context) { //发帖
	var newMainPost models.MainPost
	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	user_p := body["user_id"]
	content := body["content"]
	title := body["title"]
	db.AutoMigrate(&models.User{})
	result := db.Last(&user, "ID = ?", user_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(612)
		return
	}
	db.AutoMigrate(&models.MainPost{})
	db.Create(&models.MainPost{User_p: user_p, Title: title, Video_p: "indepent", Views: 0, ContentShow: content})
	db.Last(&newMainPost)
	mainpost_p := newMainPost.ID
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
	db.AutoMigrate(&models.User{})
	result := db.Last(&user, "ID = ?", user_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(612)
		return
	}
	db.AutoMigrate(&models.MainPost{})
	result = db.Last(&mainPost, "ID = ?", mainPost_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(622)
		return
	}
	db.AutoMigrate(&models.UnitPost{})
	db.Last(&lastPost, "mainpost_p = ?", mainPost_p)
	level := lastPost.Level + 1
	db.Create(&models.UnitPost{MainPost_p: mainPost_p, User_p: user_p, Level: level, Content: content})
	db.Update()
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
	result := db.Last(&user, "ID = ?", user_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(612)
		return
	}
	result = db.Last(&unitPost, "ID = ?", post_p)
	if result.RowsAffected == 0 {
		c.AbortWithStatus(622)
		return
	}
	db.Create(&models.Comment{Content: content, Post_p: post_p, User_p: user_p})
	c.AbortWithStatus(621)
	return
}
