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
	videoPath := "https://customer-m033z5x00ks6nunl.cloudflarestream.com/ea95132c15732412d22c1476fa83f27a/manifest/video.m3u8"
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

//视频弹幕部分

func AddBullet(c *gin.Context) {
	session := sessions.Default(c)
	author := session.Get("userid")
	uauthor := int(author.(int64))
	vid := c.PostForm("vid")
	uvid, err := strconv.Atoi(vid)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err) //返回400
		return
	}
	text := c.PostForm("text")
	time := c.PostForm("time")
	timeFloat, err := strconv.ParseFloat(time, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err) //返回400
		return
	}
	color := c.PostForm("color")
	position := c.PostForm("type")
	bullet := VideoBullet{Author: uauthor, Vid: uvid, Text: text, Time: timeFloat, Color: color, Force: false}
	switch position {
	case "scroll":
		bullet.Top = false
		bullet.Bottom = false
	case "top":
		bullet.Top = true
		bullet.Bottom = false
	case "bottom":
		bullet.Top = false
		bullet.Bottom = true
	default:
		c.AbortWithError(http.StatusBadRequest, err)
	}
	_, err1 := db.InsertOne(bullet)
	if err1 != nil {
		c.AbortWithError(http.StatusInternalServerError, err1)
	}
	return
}

func BrowseBullets(c *gin.Context) {
	session := sessions.Default(c)
	author := session.Get("userid")
	uauthor := int(author.(int64))
	vid := c.Param("id")
	var bullets []VideoBullet
	uvid, err := strconv.Atoi(vid)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	has, _ := db.In("vid", uvid).Count(&VideoBullet{})
	if has == 0 {
		println(uvid)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	db.In("vid", uvid).Asc("time").Find(&bullets)
	var bulletsOutput []gin.H
	//弹幕所处位置的格式
	for _, i := range bullets {
		var positionStr string
		if i.Top == true {
			positionStr = "top"
		} else if i.Bottom == true {
			positionStr = "bottom"
		} else {
			positionStr = "scroll"
		}
		//判断弹幕是否为自己发送
		var isMe bool
		if i.Author == uauthor {
			isMe = true
		} else {
			isMe = false
		}

		bulletsOutput = append(bulletsOutput,
			gin.H{
				"color": i.Color,
				"text":  i.Text,
				"time":  i.Time,
				"type":  positionStr,
				"isMe":  isMe,
				"force": i.Force,
			})
	}
	c.JSONP(http.StatusOK, gin.H{
		"items": bulletsOutput,
	})
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
