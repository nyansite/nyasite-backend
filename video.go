package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func NewTag(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("is_login") != true {
		c.AbortWithStatus(http.StatusUnauthorized) //返回401
		return
	}
	level := session.Get("level").(uint)
	privilege_level := level >> 4
	if privilege_level <= 10 { //10以上权限能新建tag
		c.AbortWithStatus(http.StatusForbidden) //403
		return
	}
	tagname := c.PostForm("tagname")
	if db.Take(&Tag{}, "Text = ?", tagname).RowsAffected != 0 {
		c.AbortWithStatus(StatusRepeatTag)
		return
	}
	db.Create(&Tag{Text: tagname})
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

func GetVideoImg(c *gin.Context) {

	strid := c.Param("id")
	if strid == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	id, err := strconv.Atoi(strid)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //返回400
		return
	}
	fimg, err := GetFile("/img/" + strconv.Itoa(id))
	if err != nil {
		if err == NotFound {
			c.String(http.StatusNotFound, "梦里啥都有,'"+("/img/"+strconv.Itoa(id))+"'什么都没有")
			return
		}
		c.String(http.StatusInternalServerError, err.Error()) //500
		return
	}

	c.Header("Content-Encoding", "br")  //声明压缩格式,否则会被当作二进制文件下载
	c.Header("Vary", "Accept-Encoding") //客户端使用缓存
	c.Data(http.StatusOK, "image/webp", fimg)
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

// 不审核直接上传,测试接口
func UploadVideo(c *gin.Context) {
	title := c.PostForm("Title")
	description := c.PostForm("Description") //简介
	f, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	ff, _ := f.Open()
	defer ff.Close()
	// c.SaveUploadedFile(f, )
	var video Video
	video.Title = title
	video.Description = description
	video.Views = 0
	db.Create(&video)
	err = ffmpeg.Input("").Run()
	AddFile(ff, "/video/"+strconv.Itoa(int(video.ID))+".m3u8")
}
