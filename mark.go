package main

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//返回收藏视频

func GetMarksRecord(c *gin.Context) {
	uid := GetUserIdWithoutCheck(c)
	pageStr := c.Param("pg")
	pg, _ := strconv.Atoi(pageStr)
	if pg < 1 {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	count, err := db.In("uid", uid).Count(&VideoMarkRecord{})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if count < 1 {
		c.AbortWithStatus(http.StatusNotFound)
	}
	if pg > int(math.Ceil(float64(count)/20)) {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	var marks []VideoMarkRecord
	var videosReturn []VideoReturn
	db.In("uid", uid).Desc("created_at").Find(&marks)
	for _, i := range marks {
		var video Video
		var videoReturn VideoReturn
		author := DBGetCircleDataShow(i.Uid)
		db.ID(i.Vid).Get(&video)
		videoReturn.Id = video.Id
		videoReturn.Title = video.Title
		videoReturn.CoverPath = video.CoverPath
		videoReturn.Views = video.Views - 1
		videoReturn.Likes = video.Likes - 1
		videoReturn.Author.Id = author.Id
		videoReturn.Author.Name = author.Name
		videoReturn.CreatedAt = video.CreatedAt
		videosReturn = append(videosReturn, videoReturn)
	}
	c.JSON(http.StatusOK, gin.H{
		"count": count,
		"body":  videosReturn,
	})
}
