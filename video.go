package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func NewTag(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("is_login") != true {
		c.AbortWithStatus(http.StatusUnauthorized) //返回401
		return
	}
	tagname := c.PostForm("tagname")
	if has, _ := db.Exist(&Tag{Text: tagname}); has == true {
		c.AbortWithStatus(StatusRepeatTag)
		return
	}
	db.Insert(&Tag{Text: tagname})
	c.AbortWithStatus(http.StatusOK)
}

func GetVideoComment(c *gin.Context) {
	strid := c.Param("id")
	spg := c.Param("pg")
	sid, err := strconv.Atoi(strid)
	id := uint(sid) //视频id
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //返回400
		return
	}
	pg, err := strconv.Atoi(spg)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //返回400
		return
	}
	var comments []VideoComment

	db.Desc("Likes").Limit(20, (pg-1)*20).Where("Vid = ?", id).Find(&comments)
	fmt.Println(&comments)
}

func GetVideoImg(c *gin.Context) {

	strid := c.Param("id")
	if strid == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	id, err := strconv.Atoi(strid)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //返回400 		return
	}
	var video Video
	_, err = db.ID(id).Get(&video)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error()) //500
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=cover.webp")
	c.Header("Vary", "Accept-Encoding")
	//c.Header("Content-Encoding", "br") //声明压缩格式,否则会被当作二进制文件下载
	c.File(video.CoverPath)
	return
}

// TODO 先摸了
func AddComment(c *gin.Context) {
	// var video Video
	// db.Preload("CommentP").First(&video, vid)

	// // fmt.Println(video.CommentP[len(video.CommentP)-1].ID)
	// var com []VideoComment
	// db.Find(&com, "Pid = ?", video.CommentP[len(video.CommentP)-1].ID)

	// if len(com) >= 16 {
	// 	video.CommentP = append(video.CommentP, VideoCommentPage{Vid: video.ID})

	// }
	// pg := video.CommentP[len(video.CommentP)-1]

	// if cid == 0 {
	// 	pg.Comment = append(pg.Comment, VideoComment{Text: str})
	// 	db.Save(&pg)
	// }
}

func SaveVideo(author uint, src string, cscr string, title string, description string, uuid string) {
	var video Video
	video.Author = author
	video.Title = title
	video.Description = description
	video.Views = 0
	video.CoverPath = cscr
	err := ffmpeg.Input(src).Output(src+".mp4", ffmpeg.KwArgs{
		// "c:v": "libsvtav1",
	}).Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err != nil {
		fmt.Println(err)
	}
	video.IpfsHash = Upload(src + ".mp4")
	db.Insert(&video)
	return
}
