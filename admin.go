package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	UUID "github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func AdminVideo(ctx *gin.Context) {

	vpg := ctx.Param("page")
	if vpg == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, "1")
		return
	}

	pg, err := strconv.Atoi(vpg)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	if pg < 1 {
		ctx.String(http.StatusBadRequest, "你搁这翻空气呢?")
		return
	}

	ctx.HTML(http.StatusOK, "browsevideo.html", gin.H{})
}

func AdminVideoPost(ctx *gin.Context) {
	vpg := ctx.Param("page")
	pg, err := strconv.Atoi(vpg)
	if err != nil || pg < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	var videos []Video
	var count int64 //总数,Count比rowsaffected更快(懒得用变量缓存了
	pg -= 1
	db.Model(&Video{}).Count(&count)
	db.Limit(20).Offset(pg * 20).Find(&videos)
	ctx.JSON(http.StatusOK, gin.H{
		"Body":      videos,
		"PageCount": count,
	})
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

	uuid := UUID.New()
	fpath := "./temporary/" + uuid.String() + path.Ext(f.Filename)
	c.SaveUploadedFile(f, fpath)
	go SaveVideo(fpath, title, description, uuid.String())

}

func SaveVideo(src string, title string, description string, uuid string) {
	var video Video
	video.Title = title
	video.Description = description
	video.Views = 0
	db.Create(&video)
	err := ffmpeg.Input(src).Output(src+".mp4", ffmpeg.KwArgs{"c:v": "libsvtav1"}).Run() //cpu软解
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
	}).Run() //cpu软解
	err = Addpath("./temporary/"+uuid, "/video/"+strconv.Itoa(int(video.ID)))
	if err != nil {
		fmt.Println(err)
	}
}
