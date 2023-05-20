package upload

// import (
// 	"cute_site/models"

// 	"fmt"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// func Upload(db *gorm.DB, c *gin.Context) {
// 	//获取标题和简介
// 	up_p := c.PostForm("up_p")
// 	title := c.PostForm("title")
// 	introduction := c.PostForm("introduction")
// 	fmt.Println(up_p)
// 	//新建一个待检视频的记录
// 	var videoUpload models.VideoRequireReview
// 	db.AutoMigrate(&models.VideoRequireReview{})
// 	db.Create(&models.VideoRequireReview{Title: title, Introduction: introduction, UpId_p: up_p, Pass: 0})
// 	db.Last(&videoUpload)
// 	file, err := c.FormFile("file")
// 	//将上传路径定义为/media/vntc/+待检视频的id
// 	id := fmt.Sprintf("%d", videoUpload.ID)
// 	dst := "./media/vntc/" + id + "/" + file.Filename
// 	videoUpload.VideoFile = dst
// 	//一个规范的纠错部分
// 	if err != nil {
// 		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
// 		return
// 	}
// 	if err := c.SaveUploadedFile(file, dst); err != nil {
// 		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
// 		return
// 	}
// 	cover, err := c.FormFile("cover")
// 	dst = "./media/vntc/" + id + "/" + cover.Filename
// 	videoUpload.CoverFile = dst
// 	if err != nil {
// 		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
// 		return
// 	}
// 	if err := c.SaveUploadedFile(cover, dst); err != nil {
// 		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
// 		return
// 	}
// 	db.Save(&videoUpload)
// }
