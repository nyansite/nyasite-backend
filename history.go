package main

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetHistoryRecord(c *gin.Context) {
	uid := GetUserIdWithoutCheck(c)
	pageStr := c.Param("pg")
	pg, _ := strconv.Atoi(pageStr)
	if pg < 1 {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	count, err := db.In("uid", uid).Count(&VideoPlayedRecord{})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if count < 1 {
		c.AbortWithStatus(http.StatusNotFound)
	}
	if pg > int(math.Ceil(float64(count)/20)) {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	var history []VideoPlayedRecord
	var videosReturn []VideoReturn
	db.In("uid", uid).Desc("last_play").Find(&history)
	for _, i := range history {
		var video Video
		var videoReturn VideoReturn
		author := DBGetCircleDataShow(i.Uid)
		db.ID(i.Vid).Get(&video)
		videoReturn.Id = video.Id
		videoReturn.Title = video.Title
		videoReturn.CoverPath = video.CoverPath
		videoReturn.Views = video.Views
		videoReturn.Likes = video.Likes
		videoReturn.Author.Id = author.Id
		videoReturn.Author.Name = author.Name
		//用createdAt表示上次观看的时间
		videoReturn.CreatedAt = i.LastPlay
		videosReturn = append(videosReturn, videoReturn)
	}
	c.JSON(http.StatusOK, gin.H{
		"count": count,
		"body":  videosReturn,
	})
}
