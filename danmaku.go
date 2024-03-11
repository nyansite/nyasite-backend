package main

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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
}

func BrowseBullets(c *gin.Context) {
	is_login, _ := c.Cookie("is_login")
	author := GetUserIdWithCheck(c)
	vid := c.Param("id")
	var bullets []VideoBullet
	uvid, err := strconv.Atoi(vid)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	has, _ := db.In("vid", uvid).Count(&VideoBullet{})
	if has == 0 {
		c.JSONP(http.StatusOK, gin.H{
			"autoInsert": is_login == "true",
		})
		return
	}
	db.In("vid", uvid).Asc("time").Find(&bullets)
	var bulletsOutput []gin.H
	//弹幕所处位置的格式
	for _, i := range bullets {
		var positionStr string
		if i.Top {
			positionStr = "top"
		} else if i.Bottom{
			positionStr = "bottom"
		} else {
			positionStr = "scroll"
		}
		//判断弹幕是否为自己发送
		var isMe bool
		if i.Author == author {
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
		"items":      bulletsOutput,
		"autoInsert": is_login == "true",
	})
}
