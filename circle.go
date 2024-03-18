package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用于用户界面
// 获取已经关注的社团
func GetCirclesSubscribed(c *gin.Context) {
	userid := GetUserIdWithoutCheck(c)
	var membersOfCircle []MemberOfCircle
	var circles []CircleDataShow
	has, _ := db.Where("permission = 0 and uid = ?", userid).Count(&MemberOfCircle{})
	if has == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	db.Where("permission = 0 and uid = ?", userid).Find(&membersOfCircle)
	for _, i := range membersOfCircle {
		circles = append(circles, DBGetCircleDataShow(i.Cid))
	}
	c.JSON(http.StatusOK, gin.H{
		"circles": circles,
	})
}

// 获取已经加入的社团
func GetCircleJoined(c *gin.Context) {
	userid := GetUserIdWithoutCheck(c)
	var membersOfCircle []MemberOfCircle
	var circles []CircleDataShow
	has, _ := db.Where("permission > 0 and uid = ?", userid).Count(&MemberOfCircle{})
	if has == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	db.Where("permission > 0 and uid = ?", userid).Find(&membersOfCircle)
	for _, i := range membersOfCircle {
		circle := DBGetCircleDataShow(i.Cid)
		circles = append(circles, circle)
	}
	c.JSON(http.StatusOK, gin.H{
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
	if kind >= 3 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	userid := GetUserIdWithoutCheck(c)
	var authorOfCircle []MemberOfCircle
	var circles []gin.H
	has, _ := db.In("uid", userid).Count(&MemberOfCircle{})
	if has == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.Where("uid = ? and permission >= 2", userid).Find(&authorOfCircle)
	var circle Circle
	for _, i := range authorOfCircle {
		db.ID(i.Cid).Get(&circle)
		if (circle.Kinds & int16(1<<kind)) > 0 { //压位
			circles = append(circles, gin.H{
				"name": circle.Name,
				"id":   i.Cid,
			})
		}
	}
	if len(circles) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, gin.H{
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
	_, err := db.InsertOne(&circleApplication)
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
		db.Where("cid = ? and uid = ?", vCid, uid).Unscoped().Delete(&MemberOfCircle{})
		return
	}
	memberOfCircle := MemberOfCircle{
		Cid:        vCid,
		Uid:        uid,
		Permission: 0,
	}
	db.InsertOne(&memberOfCircle)
}

func GetCircle(c *gin.Context) {
	strCid := c.Param("id")
	vCid, _ := strconv.Atoi(strCid)
	var circle Circle
	exist, _ := db.ID(vCid).Get(&circle)
	if exist == false {
		c.AbortWithStatus(http.StatusNotFound)
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
	c.JSON(http.StatusOK, gin.H{
		"id":         circle.Id,
		"name":       circle.Name,
		"avatar":     circle.Avatar,
		"descrption": circle.Descrption,
		"kinds":      circle.Kinds,
		"members":    DBGetMembersOfCircle(int(circle.Id)),
		"relation":   DBgetRelationToCircle(int(circle.Id), c),
		"videos":     videosDisplay,
		"createdAt":  circle.CreatedAt,
	})
}

func GetVideosOfCircle(c *gin.Context) {
	strId := c.Param("id")
	strPage := c.Param("page")
	strMethod := c.Param("method")
	id, _ := strconv.Atoi(strId)
	page, _ := strconv.Atoi(strPage)
	method, err := strconv.Atoi(strMethod)
	if err != nil {
		method = 0
	}
	var videos []Video
	var videosReturn []VideoReturn
	var count int64
	switch method {
	case 0:
		count, _ = db.In("author", id).Limit(20, (page-1)*20).Desc("id").FindAndCount(&videos)
	case 1:
		count, _ = db.In("author", id).Limit(20, (page-1)*20).Asc("views").FindAndCount(&videos)
	default:
		c.AbortWithStatus(http.StatusBadRequest)
	}
	for _, i := range videos {
		var videoReturn VideoReturn
		videoReturn.CoverPath = i.CoverPath
		videoReturn.CreatedAt = i.CreatedAt
		videoReturn.Id = i.Id
		videoReturn.Likes = i.Likes
		videoReturn.Title = i.Title
		videoReturn.Views = i.Views
		videosReturn = append(videosReturn, videoReturn)
	}
	c.JSON(http.StatusOK, gin.H{
		"count":   count,
		"content": videosReturn,
	})
}

func DBGetMembersOfCircle(cid int) []UserDataShow {
	var members []MemberOfCircle
	var membersDisplay []UserDataShow
	db.Where("permission > 0 and cid = ?", cid).Find(&members)
	for _, i := range members {
		memberDisplay := DBGetUserDataShow(i.Uid)
		membersDisplay = append(membersDisplay, memberDisplay)
	}
	return membersDisplay
}

func DBgetCirclesRelatedTo(uid int) []int {
	var circlesId []int
	var membersOfCircle []MemberOfCircle
	db.Where("permission > 0 and uid = ?", uid).Find(&membersOfCircle)
	for _, i := range membersOfCircle {
		circlesId = append(circlesId, i.Cid)
	}
	return circlesId
}

func DBgetRelationToCircle(cid int, c *gin.Context) int8 {
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
