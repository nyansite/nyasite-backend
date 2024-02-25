package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用于用户界面(获取已经加入的社团)
func GetCircleJoined(c *gin.Context) {
	uauthor := GetUserIdWithoutCheck(c)
	var membersOfCircle []MemberOfCircle
	var circles []Circle
	has, _ := db.Where("permission > 0 and uid = ?", uauthor).Count(&MemberOfCircle{})
	if has == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	db.Where("permission > 0 and uid = ?", uauthor).Find(&membersOfCircle)
	for _, i := range membersOfCircle {
		var circle Circle
		db.ID(i.Cid).Get(&circle)
		circles = append(circles, circle)
	}
	c.JSONP(http.StatusOK, gin.H{
		"circles": circles,
	})
}

// 用于上传界面
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
	uauthor := GetUserIdWithoutCheck(c)
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
		if i.Permission >= 2 && ((circle.Kinds & int16(1<<kind)) > 0) { //压位
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

func SubscribeCircle(c *gin.Context) {
	strCid := c.PostForm("cid")
	vCid, _ := strconv.Atoi(strCid)
	uid := GetUserIdWithoutCheck(c)
	has, _ := db.Where("cid = ? and uid = ?", vCid, uid).Count(&MemberOfCircle{})
	if has != 0 {
		db.Where("cid = ? and uid = ?", vCid, uid).Delete(&MemberOfCircle{})
		return
	}
	memberOfCircle := MemberOfCircle{
		Cid:        vCid,
		Uid:        uid,
		Permission: 0,
	}
	db.Insert(&memberOfCircle)
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
		"id":         circle.Id,
		"name":       circle.Name,
		"avatar":     circle.Avatar,
		"descrption": circle.Descrption,
		"kinds":      circle.Kinds,
		"members":    membersDisplay,
		"relation":   DBgetRelationWithCircle(int(circle.Id), c),
		"videos":     videosDisplay,
		"createdAt":  circle.CreatedAt,
	})
}

func DBgetRelationWithCircle(cid int, c *gin.Context) int8 {
	uid := GetUserIdWithCheck(c)
	if uid == -1 {
		return -2
	}
	var memberOfCircle MemberOfCircle
	has, _ := db.Where("cid = ? and uid = ?", cid, uid).Count(&MemberOfCircle{})
	if has == 0 {
		return -1
	}
	db.Where("cid = ? and uid = ?", cid, uid).Get(&memberOfCircle)
	return int8(memberOfCircle.Permission)
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
