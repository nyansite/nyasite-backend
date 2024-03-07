package main

import (
	"fmt"
	"net/http"
	"strconv"

	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAllCirclesNeedtoCheck(c *gin.Context) {
	uauthor := GetUserIdWithoutCheck(c)
	var circlesNeedtoCheck []ApplyCircle
	db.In("stauts", false).Find(&circlesNeedtoCheck)
	var circlesNeedtoCheckDisplay []ApplyCircle
	var count int64
	for _, i := range circlesNeedtoCheck {
		count, _ = db.Where("acid = ? AND reviewer = ?", i.Id, uauthor).Asc("Id").Count(&VoteOfApplyCircle{})
		if count == 0 {
			circlesNeedtoCheckDisplay = append(circlesNeedtoCheckDisplay, i)
		}
	}
	c.JSON(http.StatusOK, gin.H{"results": circlesNeedtoCheckDisplay})
}

func VoteForCirclesNeedtoCheck(c *gin.Context) {
	uauthor := GetUserIdWithoutCheck(c)
	acid := c.PostForm("acid")
	acidNumber, _ := strconv.Atoi(acid)
	altitude := c.PostForm("altitude")
	var altitudeBool bool
	altitudeBool, err5 := strconv.ParseBool(altitude)
	if err5 != nil {
		c.AbortWithError(http.StatusBadRequest, err5)
	}
	agree, _ := db.Where("acid = ? AND agree = ?", acid, true).Count(&VoteOfApplyCircle{})
	disagree, _ := db.Where("acid = ? AND agree = ?", acid, false).Count(&VoteOfApplyCircle{})
	var circleNeedToCheck ApplyCircle
	_, err4 := db.ID(acidNumber).Get(&circleNeedToCheck)
	if err4 != nil {
		c.AbortWithError(http.StatusInternalServerError, err4)
	}
	if agree >= 4 && altitudeBool {
		circle := Circle{
			Name:       circleNeedToCheck.Name,
			Avatar:     circleNeedToCheck.Avatar,
			Descrption: circleNeedToCheck.Descrption,
			Kinds:      circleNeedToCheck.Kinds,
		}
		_, err := db.Insert(&circle)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		memberOfCircle := MemberOfCircle{
			Uid:        uauthor,
			Cid:        int(circle.Id),
			Permission: 4,
		}
		_, err2 := db.Insert(&memberOfCircle)
		if err2 != nil {
			c.AbortWithError(http.StatusInternalServerError, err2)
			return
		}
		_, err1 := db.ID(acidNumber).Unscoped().Delete(&ApplyCircle{})
		if err1 != nil {
			c.AbortWithError(http.StatusInternalServerError, err1)
			return
		}
		_, err3 := db.In("acid", acidNumber).Delete(&VoteOfApplyCircle{})
		if err3 != nil {
			c.AbortWithError(http.StatusInternalServerError, err3)
			return
		}
		return
	} else if disagree >= 4 && !altitudeBool {
		circleNeedToCheck.Stauts = true
		_, err1 := db.ID("acid").Cols("stauts").Update(&circleNeedToCheck)
		if err1 != nil {
			c.AbortWithError(http.StatusInternalServerError, err1)
			return
		}
	} else {
		voteOfApplyCircle := VoteOfApplyCircle{
			Reviewer: uauthor,
			Agree:    altitudeBool,
			Acid:     acidNumber,
		}
		db.Insert(&voteOfApplyCircle)
		return
	}
}

func GetAllVideoNeedToChenck(c *gin.Context) {
	var videosNeedToCheck []VideoNeedToCheck
	db.In("stauts", false).Find(&videosNeedToCheck)
	c.JSONP(http.StatusOK, gin.H{"results": videosNeedToCheck})
}

func PassVideo(c *gin.Context) {
	videoNeedToCheckId := c.PostForm("vcid")
	var videoNeedToCheck VideoNeedToCheck
	bool, err := db.ID(videoNeedToCheckId).Get(&videoNeedToCheck)
	if bool == false {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	video := Video{
		CoverPath:   videoNeedToCheck.CoverPath,
		VideoUid:    videoNeedToCheck.VideoUid,
		Title:       videoNeedToCheck.Title,
		Description: videoNeedToCheck.Description,
		Author:      videoNeedToCheck.Author,
		Views:       1,
		Likes:       1,
		Caution:     1,
		Upid:        videoNeedToCheck.Upid,
	}
	_, err1 := db.Insert(&video)
	if err1 != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	for _, i := range videoNeedToCheck.Tags {
		tag := Tag{
			Tid: i,
			Pid: int(video.Id),
		}
		fmt.Println(tag)
		db.Insert(&tag)
	}
	db.ID(videoNeedToCheckId).Delete(&VideoNeedToCheck{})
	return
}

func RejectVideo(c *gin.Context) {
	videoNeedToCheckId := c.PostForm("vcid")
	vcidNumber, _ := strconv.Atoi(videoNeedToCheckId)
	var videoNeedToCheck VideoNeedToCheck
	_, err := db.ID(vcidNumber).Get(&videoNeedToCheck)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	videoNeedToCheck.Stauts = true
	_, err1 := db.ID(vcidNumber).Cols("stauts").Update(&videoNeedToCheck)
	if err1 != nil {
		c.AbortWithError(http.StatusInternalServerError, err1)
	}
	return
}
