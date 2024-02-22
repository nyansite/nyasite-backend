package main

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CheckAvailableCircle(c *gin.Context) {
	strKind := c.Param("type")
	kind, err := strconv.Atoi(strKind)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	if kind > 3 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	session := sessions.Default(c)
	author := session.Get("userid")
	uauthor := int(author.(int64))
	var authorOfCircle []MemberOfCircle
	var circles []gin.H
	has, _ := db.In("uid", uauthor).Count(&MemberOfCircle{})
	if has == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.In("uid", uauthor).Asc("permission").Find(&authorOfCircle)
	var circle Circle
	for _, i := range authorOfCircle {
		db.ID(i.Cid).Get(&circle)
		if i.Permission <= 2 && ((circle.Kinds & int16(1<<kind)) > 0) { //压位
			circles = append(circles, gin.H{
				"name": circle.Name,
				"id":   i.Cid,
			})
		}
	}
	if len(circles) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSONP(http.StatusOK, gin.H{
			"circles": circles,
		})
	}
}

func PostCircleApplication(c *gin.Context) {
	author := GetUserIdWithoutCheck(c)
	name := c.PostForm("name")
	pastworks := c.PostForm("pastworks")
	kindsStr := c.PostFormArray("type")
	avatar := c.PostForm("avatar")
	description := c.PostForm("description")
	var kinds int16
	for _, i := range kindsStr {
		switch i {
		case "video":
			kinds = kinds + (1 << 0)
		case "music":
			kinds = kinds + (1 << 1)
		case "image":
			kinds = kinds + (1 << 2)
		}
	}
	has, _ := db.In("name", name).Count(&Circle{})
	if has > 0 {
		c.AbortWithStatus(http.StatusInsufficientStorage)
		return
	}
	circleApplication := ApplyCircle{
		Name:       name,
		Avatar:     avatar,
		Descrption: description,
		ApplyText:  pastworks,
		Stauts:     false,
		Kinds:      kinds,
		Applicant:  author,
	}
	_, err := db.Insert(&circleApplication)
	if err != nil {
		c.AbortWithError(http.StatusInsufficientStorage, err)
		return
	}
	return
}

func DBGetCircleDataShow(cid int) CircleDataShow {
	var circle Circle
	db.ID(cid).Get(&circle)
	circleDisplay := CircleDataShow{
		Id:     circle.Id,
		Name:   circle.Name,
		Avatar: circle.Avatar,
	}
	return circleDisplay
}

func GetCircle(c *gin.Context) {
	strCid := c.Param("id")
	vCid, _ := strconv.Atoi(strCid)
	var circle Circle
	exist, _ := db.ID(vCid).Get(&circle)
	if exist == false {
		c.AbortWithStatus(http.StatusNotFound)
	}
	var members []MemberOfCircle
	var membersDisplay []UserDataShow
	db.In("cid", vCid).Find(&members)
	for _, i := range members {
		memberDisplay := DBGetUserDataShow(i.Uid)
		membersDisplay = append(membersDisplay, memberDisplay)
	}
	var videos []Video
	var videosDisplay []VideoReturn
	db.In("author", vCid).Limit(20, 0).Desc("id").Find(&videos)
	for _, i := range videos {
		var videoReturn VideoReturn
		videoReturn.CoverPath = i.CoverPath
		videoReturn.CreatedAt = i.CreatedAt
		videoReturn.Id = i.Id
		videoReturn.Title = i.Title
		videoReturn.Views = i.Views
		videosDisplay = append(videosDisplay, videoReturn)
	}
	c.JSONP(http.StatusOK, gin.H{
		"name":       circle.Name,
		"avatar":     circle.Avatar,
		"descrption": circle.Descrption,
		"kinds":      circle.Kinds,
		"members":    membersDisplay,
		"videos":     videosDisplay,
		"createdAt":  circle.CreatedAt,
	})
}
