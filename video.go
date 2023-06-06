package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os"
)

func NewTag(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("is_login") != true {
		c.AbortWithStatus(http.StatusUnauthorized) //返回401
		return
	}
	level := session.Get("level").(uint)
	privilege_level := level >> 4
	if privilege_level < 10 { //10以上权限能新建tag
		c.AbortWithStatus(http.StatusForbidden) //403
		return
	}
	tagname := c.PostForm("tagname")
	if has, _ :=db.Exist(&Tag{Text: tagname}); has == true {
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
	// db.Preload("CommentPage", "Count = ?", pg).Preload("Cid = 0 OR (Count < 3)").First(&video, id)
	// db.Limit(3).Order("likes").Preload("VideoCommentReply").Limit(20).Offset(pg-1).Where("Vid = ?", id).Find(&comments)
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

func SaveVideo(src string, title string, description string, uuid string) {
	var video Video
	video.Title = title
	video.Description = description
	video.Views = 0
	db.Insert(&video)
	err := ffmpeg.Input(src).Output(src+".mp4", ffmpeg.KwArgs{
		// "c:v": "libsvtav1",
	}).Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Mkdir("./temporary/"+uuid, os.ModePerm)
	err = ffmpeg.Input(src+".mp4").Output("./temporary/"+uuid+"/w.m3u8", ffmpeg.KwArgs{
		// "codec":         "copy",		//只有音频,原因未知
		"start_number":  0,
		"hls_list_size": 0,
		"hls_time":      5,
		"f":             "hls",
	}).Run()
	err = Addpath("./temporary/"+uuid, "/video/"+strconv.Itoa(int(video.Id)))
	if err != nil {
		fmt.Println(err)
	}

}
