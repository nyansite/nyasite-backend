package main

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
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
	c.JSON(http.StatusOK, gin.H{
		"title":       video.Title,
		"videoPath":   videoPath,
		"author":      video.Author,
		"creatTime":   video.CreatedAt,
		"description": video.Description,
		"views":       video.Views,
		"likes":       video.Likes,
	})
}

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
	session := sessions.Default(c)
	author := session.Get("userid")
	uauthor := int(author.(int64))
	title := c.PostForm("title")
	description := c.PostForm("description")
	cover := c.PostForm("cover")
	tags := c.PostFormArray("tags")
	var err error
	var tagsUint8 []uint8
	var unitTag int
	for _, i := range tags {
		unitTag, err = strconv.Atoi(i)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		tagsUint8 = append(tagsUint8, uint8(unitTag))
	}
	newVideo := VideoNeedToCheck{Author: uauthor, Title: title, Description: description, CoverPath: cover, Tags: tagsUint8}
	_, err1 := db.InsertOne(newVideo)
	if err1 != nil {
		c.AbortWithError(http.StatusInternalServerError, err1)
	}
	return
}
