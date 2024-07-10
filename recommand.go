package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RankVideos(frequency uint8) {
	//0:daily 1: every week 2.every month
	var videos []Video
	db.Find(&videos)
	var videosDataToSort []VideoReturn
	for _, i := range videos {
		var videoData VideoReturn
		videoData.Id = i.Id
		videoData.Likes = i.Likes
		videoData.Views = i.Views
		videosDataToSort = append(videosDataToSort, videoData)
	}
	SortVideo(videosDataToSort, func(p, q *VideoReturn) bool {
		countCommentsP, _ := db.In("vid", p.Id).Count(&VideoComment{})
		countCommentsQ, _ := db.In("vid", q.Id).Count(&VideoComment{})
		countBulletsP, _ := db.In("vid", p.Id).Count(&VideoBullet{})
		countBulletsQ, _ := db.In("vid", q.Id).Count(&VideoBullet{})
		//权重 = 点赞*1.2 + 浏览量 + 评论数*0.7 + 弹幕数*0.5
		pIndex := float64(p.Likes)*1.2 + float64(p.Views) + float64(countCommentsP)*0.7 + float64(countBulletsP)*0.5
		qIndex := float64(q.Likes)*1.2 + float64(q.Views) + float64(countCommentsQ)*0.7 + float64(countBulletsQ)*0.5
		return pIndex > qIndex
	})
	for _, i := range videosDataToSort[:14] {
		trendingVideo := TrendingRankVideo{Vid: int(i.Id), Type: frequency}
		db.InsertOne(&trendingVideo)
	}
	return
}

func RankVideosDaliy()   { RankVideos(0) }
func RankVideosMouthly() { RankVideos(1) }
func RankVideosYearly()  { RankVideos(2) }

func RankVideosTest(c *gin.Context) {
	frequencyStr := c.PostForm("type")
	frequency, err := strconv.Atoi(frequencyStr)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	if frequency == 3 {
		RankVideos(0)
		RankVideos(1)
		RankVideos(2)
	} else if frequency < 3 && frequency >= 0 {
		RankVideos(uint8(frequency))
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return
}

func GetTrending(c *gin.Context) {
	dailyTrendingVideos := getTrendingVideosReturn(0)
	weeklyTrendingVideos := getTrendingVideosReturn(1)
	monthlyTrendingVideos := getTrendingVideosReturn(2)
	c.JSON(http.StatusOK, gin.H{
		"daily":   dailyTrendingVideos,
		"weekly":  weeklyTrendingVideos,
		"monthly": monthlyTrendingVideos,
	})
}

func getTrendingVideosReturn(frq uint8) []VideoReturn {
	var trendingVideos []VideoReturn
	var trending []TrendingRankVideo
	db.In("type", frq).Find(&trending)
	for _, i := range trending {
		var video Video
		var videoReturn VideoReturn
		db.ID(i.Vid).Get(&video)
		author := DBGetCircleDataShow(video.Author)
		videoReturn.Id = video.Id
		videoReturn.Views = video.Views - 1
		videoReturn.Likes = video.Likes - 1
		videoReturn.Title = video.Title
		videoReturn.CoverPath = video.CoverPath
		videoReturn.Author.Id = author.Id
		videoReturn.Author.Name = author.Name
		videoReturn.CreatedAt = video.CreatedAt
		trendingVideos = append(trendingVideos, videoReturn)
	}
	return trendingVideos
}
