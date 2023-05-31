package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

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

func GetVideoComment(c *gin.Context) {
	strid := c.Param("id")
	spg := c.Param("pg")
	sid, err := strconv.Atoi(strid)
	id := uint(sid)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //返回400
		return
	}
	pg, err := strconv.Atoi(spg)
	// pg = uint(pg)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //返回400
		return
	}
	var video Video
	db.Preload("CommentPage", "Count = ?", pg).Preload("Cid = 0 OR (Count < 3)").First(&video, id)
	fmt.Println(&video)
}

func AddComment(c *gin.Context) {

}

// TODO 先摸了
func DaddComment(str string, vid uint, cid uint) {
	var video Video
	db.Preload("CommentP").First(&video, vid)

	// fmt.Println(video.CommentP[len(video.CommentP)-1].ID)
	var com []VideoComment
	db.Find(&com, "Pid = ?", video.CommentP[len(video.CommentP)-1].ID)

	if len(com) >= 16 {
		video.CommentP = append(video.CommentP, VideoCommentPage{Vid: video.ID})

	}
	pg := video.CommentP[len(video.CommentP)-1]

	if cid == 0 {
		pg.Comment = append(pg.Comment, VideoComment{Text: str})
		db.Save(&pg)
	}
}

// 先摸了
func uploadVideo(c *gin.Context) {
	//获取标题和简介
	upid := c.PostForm("upid")
	title := c.PostForm("title")
	profile := c.PostForm("profile")
	fmt.Println(upid)
	file, err := c.FormFile("file")

	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	cover, err := c.FormFile("cover")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	DuploadVideo(upid, title, profile, file, cover, c)
	return
}







func DuploadVideo(upid string, title string, profile string,
	file *multipart.FileHeader, cover *multipart.FileHeader, c *gin.Context) {

	nid, _ := strconv.Atoi(upid)
	videoUpload := VideoPreviewRequire{Title: title, Profile: profile, Up: uint(nid), Pass: 0}
	db.Create(&videoUpload)
	//将上传路径定义为/media/vntc/+待检视频的id
	id := fmt.Sprintf("%d", videoUpload.ID)
	dst := "./media/vntc/" + id + "/" + file.Filename
	videoUpload.VideoFile = dst
	//一个规范的纠错部分
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}
	dst = "./media/vntc/" + id + "/" + cover.Filename
	videoUpload.CoverFile = dst
	if err := c.SaveUploadedFile(cover, dst); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}
	db.Save(&videoUpload)
	return
}
