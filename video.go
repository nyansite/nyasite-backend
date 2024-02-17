package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//视频返回

func GetVideo(c *gin.Context) {
	strVid := c.Param("id")
	vVid, _ := strconv.Atoi(strVid)
	var video Video
	exist, _ := db.ID(int(vVid)).Get(&video)
	if exist == false {
		c.AbortWithStatus(http.StatusNotFound)
	}
	//获取路径
	//test data
	videoPath := "https://customer-f33zs165nr7gyfy4.cloudflarestream.com/6b9e68b07dfee8cc2d116e4c51d6a957/manifest/video.m3u8"
	//获取作者
	author := DBGetCircleDataShow(int(video.Id))
	c.JSON(http.StatusOK, gin.H{
		"title":       video.Title,
		"videoPath":   videoPath,
		"author":      author,
		"creatTime":   video.CreatedAt,
		"description": video.Description,
		"views":       video.Views,
		"likes":       video.Likes,
	})
}

//添加视频标签

func AddVideoTag(c *gin.Context) {
	strVid := c.PostForm("vid")
	vVid, _ := strconv.Atoi(strVid)
	uVid := int(vVid)
	strTagId := c.PostForm("tagid")
	vTagId, _ := strconv.Atoi(strTagId)
	uTagId := int(vTagId)
	tag := Tag{Tid: uTagId, Pid: uVid}
	db.Insert(tag)
	return
}

//上传视频

func PostVideo(c *gin.Context) {
	uauthor := GetUserIdWithoutCheck(c)
	title := c.PostForm("title")
	description := c.PostForm("description")
	cover := c.PostForm("cover")
	strTags := c.PostFormArray("tags")
	var err error
	var tags []int
	var tag int
	for _, i := range strTags {
		tag, err = strconv.Atoi(i)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		tags = append(tags, tag)
	}
	newVideo := VideoNeedToCheck{Author: uauthor, Title: title, Description: description, CoverPath: cover, Tags: tags, Stauts: false}
	_, err1 := db.InsertOne(newVideo)
	if err1 != nil {
		c.AbortWithError(http.StatusInternalServerError, err1)
	}
	return
}

//获取标签

func GetVideoTags(c *gin.Context) {
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
	var tagsDisplay []gin.H
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
		tagsDisplay = append(tagsDisplay, gin.H{
			"id":   tid,
			"text": tagModel.Text,
		})
	}
	c.JSONP(http.StatusOK, gin.H{
		"tags": tagsDisplay,
	})
	return
}

//获取所有视频

func GetAllVideos(c *gin.Context) {
	var videos []Video
	var videosDisplay []gin.H
	err := db.Desc("id").Find(&videos)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var author UserDataShow
	for _, i := range videos {
		author = DBGetUserDataShow(i.Author)
		videosDisplay = append(videosDisplay, gin.H{
			"id":     i.Id,
			"author": author,
			"cover":  i.CoverPath,
		})
	}
}
