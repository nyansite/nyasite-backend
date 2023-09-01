package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func AddVideoTag(c *gin.Context) {
	strVid := c.PostForm("vid")
	vVid, _ := strconv.Atoi(strVid)
	uVid := int(vVid)
	strTagId := c.PostForm("tagid")
	vTagId, _ := strconv.Atoi(strTagId)
	uTagId := int(vTagId)
	DBaddVideoTag(uVid, uTagId)
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
	return
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

func GetVideoTags(c *gin.Context) {
	var tagTexts []string
	var tagIds []int
	strid := c.Param("id")
	if strid == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	id, err := strconv.Atoi(strid)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //返回400 		return
	}
	var tags []Tag
	count, _ := db.Where("kind = ? AND pid = ?", 1, id).Count(&tags)
	if count == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	db.Where("kind = ? AND pid = ?", 1, id).Find(&tags)
	var tagModel TagModel
	var tid int
	for _, value := range tags {
		tid = int(value.Tid)
		db.ID(tid).Get(&tagModel)
		tagTexts = append(tagTexts, tagModel.Text)
		tagIds = append(tagIds, tid)
	}
	c.JSONP(http.StatusOK, gin.H{
		"tagtext": tagTexts,
		"tagid":   tagIds,
	})
	return
}

// TODO 先摸了
func AddVideoComment(ctx *gin.Context) {
	session := sessions.Default(ctx)
	author := session.Get("userid")
	vauthor := author.(int64)
	uauthor := int(vauthor)
	vid, text := ctx.PostForm("vid"), ctx.PostForm("text")
	vvid, _ := strconv.Atoi(vid)
	uvid := int(vvid)
	DBaddVideoComment(uvid, uauthor, text)
	return
}

func SaveVideo(author int, src string, cscr string, title string, description string, uuid string) {
	var video Video
	video.Author = author
	video.Title = title
	video.Description = description
	video.Views = 0
	video.CoverPath = cscr
	//the error ffmpeg part
	err := ffmpeg.Input(src).Output(src+".mp4", ffmpeg.KwArgs{
		// "c:v": "libsvtav1",
	}).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		panic(err)
	}
	//
	video.IpfsHash = Upload(src + ".mp4")
	db.Insert(&video)
	return
}

func DBaddVideoComment(vid int, author int, text string) {
	vComment := VideoComment{Vid: vid, Author: author, Text: text, Likes: 0}
	db.Insert(vComment)
	return
}

func DBaddVideoTag(vid int, tagid int) {
	tag := Tag{Tid: tagid, Kind: 1, Pid: vid}
	db.Insert(tag)
	return
}
