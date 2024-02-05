package main

import (
	"net/http"
	//"strconv"

	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAllVNTC(c *gin.Context) {

}

func PassVideo(c *gin.Context) {
	videoNeedToCheckId := c.PostForm("vid")
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
	var tag Tag
	for _, i := range videoNeedToCheck.Tags {
		tag.Tid = i
		tag.Pid = int(video.Id)
		db.Insert(tag)
	}
	db.ID(videoNeedToCheckId).Delete(&VideoNeedToCheck{})
	return
}
