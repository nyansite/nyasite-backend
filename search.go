package main

import (
	"net/http"

	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2" //看看文档，这里大量用到set特性
	"github.com/gin-gonic/gin"
)

func EnireTag(c *gin.Context) {
	var tagMs []TagModel
	db.Find(&tagMs)
	c.JSONP(http.StatusOK, gin.H{"results": tagMs})
}

func SearchVideos(c *gin.Context) {
	var videosReturn []VideoReturn
	var videoReturn VideoReturn
	vidsCount := make(map[uint]uint)
	vids := mapset.NewSet[uint]()
	tagCondition := c.PostForm("tags")
	textCondition := c.PostForm("text")
	page := c.PostForm("page")
	kind := c.PostForm("kind")
	vPage, _ := strconv.Atoi(page)
	if tagCondition == ";" {
		if textCondition == "" {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			var videos []Video
			db.Where("title like ? or description like ?", "%"+textCondition+"%", "%"+textCondition+"%").Find(&videos)
			for _, i := range videos {
				author := DBGetCircleDataShow(i.Author)
				videoReturn.Id = i.Id
				videoReturn.Title = i.Title
				videoReturn.CoverPath = i.CoverPath
				videoReturn.Views = i.Views
				videoReturn.Likes = i.Likes
				videoReturn.Author.Id = author.Id
				videoReturn.Author.Name = author.Name
				videoReturn.CreatedAt = i.CreatedAt
				videosReturn = append(videosReturn, videoReturn)
			}
		}
	} else {
		sTags := strings.Split(tagCondition, ";") //将用";"链接的tags字符串转化为slice
		sTags = sTags[:len(sTags)-1]
		nTags := uint(len(sTags)) //统计tag种类的总数
		for _, i := range sTags {
			var tagModel TagModel
			var tags []Tag
			db.In("text", i).Get(&tagModel)
			db.In("tid", &tagModel.Id).Find(&tags)
			for _, j := range tags {
				vids.Add(uint(j.Pid))
				vidsCount[uint(j.Pid)]++ //如果视频在tags出现一次就+1
			}
		}
		//bucket sort
		for i := range vids.Iter() {
			if vidsCount[i] == nTags { //如果视频出现在每一种tag中，vidsCount对应项与tag种类总数相等
				var video Video
				db.ID(i).Get(&video)
				if strings.Contains(video.Title, textCondition) || strings.Contains(video.Description, textCondition) || textCondition == "" {
					//判断是否又对应text部分,如果没有text部分，就忽略text部分
					author := DBGetCircleDataShow(video.Author)
					videoReturn.Id = video.Id
					videoReturn.Title = video.Title
					videoReturn.CoverPath = video.CoverPath
					videoReturn.Views = video.Views
					videoReturn.Likes = video.Likes
					videoReturn.Author.Id = author.Id
					videoReturn.Author.Name = author.Name
					videoReturn.CreatedAt = video.CreatedAt
					videosReturn = append(videosReturn, videoReturn)
				}
			}
		}
	}
	switch kind {
	case "0":
		SortVideo(videosReturn, func(p, q *VideoReturn) bool {
			return p.Id > q.Id
		})
	case "1":
		SortVideo(videosReturn, func(p, q *VideoReturn) bool {
			return p.Likes > q.Likes
		})
	case "2":
		SortVideo(videosReturn, func(p, q *VideoReturn) bool {
			return p.Views > q.Views
		})
	default:
		c.AbortWithStatus(http.StatusBadRequest)
	}
	var upper int
	if (vPage*20 - 1) <= len(videosReturn) {
		upper = (vPage * 20)
	} else {
		upper = len(videosReturn)
	}
	count := len(videosReturn)
	videosReturn = videosReturn[(vPage-1)*20 : upper]
	if count == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"videos": videosReturn,
		"count":  count,
	})
}
