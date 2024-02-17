package main

import (
	"net/http"
	"sort"

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

func SearchTag(c *gin.Context) {
	var tagMs []TagModel
	var tagList []string
	db.Find(&tagMs)
	for _, i := range tagMs {
		tagList = append(tagList, i.Text)
	}
	c.JSON(http.StatusOK, gin.H{"results": tagList})
	return
}

type VideoWrapper struct {
	videos []SearchVideoReturn
	by     func(p, q *SearchVideoReturn) bool
}

type SortBy func(p, q *SearchVideoReturn) bool

func (vw VideoWrapper) Len() int {
	return len(vw.videos)
}

func (vw VideoWrapper) Swap(i, j int) {
	vw.videos[i], vw.videos[j] = vw.videos[j], vw.videos[i]
}

func (vw VideoWrapper) Less(i, j int) bool {
	return vw.by(&vw.videos[i], &vw.videos[j])
}

func SortVideo(videos []SearchVideoReturn, by SortBy) {
	sort.Sort(VideoWrapper{videos, by})
}

func SearchVideos(c *gin.Context) {
	var videosReturn []SearchVideoReturn
	vidsCount := make(map[uint]uint)
	vids := mapset.NewSet[uint]()
	tagCondition := c.PostForm("tags")
	textCondition := c.PostForm("text")
	page := c.PostForm("page")
	kind := c.PostForm("kind")
	vPage, _ := strconv.Atoi(page)
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
	var author Circle
	for i := range vids.Iter() {
		if vidsCount[i] == nTags { //如果视频出现在每一种tag中，vidsCount对应项与tag种类总数相等
			var video Video
			var videoReturn SearchVideoReturn
			db.ID(i).Get(&video)
			if strings.Contains(video.Title, textCondition) || strings.Contains(video.Description, textCondition) || tagCondition == "" {
				//判断是否又对应text部分,如果没有text部分，就忽略text部分
				db.ID(video.Author).Get(&author)
				videoReturn.Id = video.Id
				videoReturn.Title = video.Title
				videoReturn.CoverPath = video.CoverPath
				videoReturn.Views = video.Views
				videoReturn.Likes = video.Likes
				videoReturn.Author.Id = author.Id
				videoReturn.Author.Name = author.Name
				videosReturn = append(videosReturn, videoReturn)
			}
		}
	}
	switch kind {
	case "0":
		SortVideo(videosReturn, func(p, q *SearchVideoReturn) bool {
			return p.Id > q.Id
		})
	case "1":
		SortVideo(videosReturn, func(p, q *SearchVideoReturn) bool {
			return p.Likes > q.Likes
		})
	case "2":
		SortVideo(videosReturn, func(p, q *SearchVideoReturn) bool {
			return p.Views > q.Views
		})
	default:
		c.AbortWithStatus(http.StatusBadRequest)
	}
	var upper int
	if (vPage*20 - 1) <= len(videosReturn) {
		upper = (vPage*20 - 1)
	} else {
		upper = len(videosReturn)
	}
	videosReturn = videosReturn[(vPage-1)*20 : upper]
	count := len(videosReturn)
	c.JSON(http.StatusOK, gin.H{
		"videos": videosReturn,
		"count":  count,
	})
}
