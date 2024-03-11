package main

import (
	"fmt"
	"net/http"
	"strconv"

	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)
func GetAllVideoNeedToChenck(c *gin.Context) {
	var videosNeedToCheck []VideoNeedToCheck
	db.In("stauts", false).Find(&videosNeedToCheck)
	c.JSONP(http.StatusOK, gin.H{"results": videosNeedToCheck})
}

func PassVideo(c *gin.Context) {
	videoNeedToCheckId := c.PostForm("vcid")
	var videoNeedToCheck VideoNeedToCheck
	bool, err := db.ID(videoNeedToCheckId).Get(&videoNeedToCheck)
	if !bool{
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
}
