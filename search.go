package main

import (
	"net/http"
	"strings"

	//"fmt"

	mapset "github.com/deckarep/golang-set/v2" //看看文档，这里大量用到set特性
	"github.com/gin-gonic/gin"
)

func EnireTag(c *gin.Context) {
	var tagMs []TagModel
	db.Find(&tagMs)
	c.JSON(http.StatusOK, gin.H{"results": tagMs})
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

func SearchVideos(c *gin.Context) {
	var tagModel TagModel
	var tags []Tag
	var videosReturn []SearchVideoReturn
	var videoReturn SearchVideoReturn
	var video Video
	vidsCount := make(map[uint]uint)
	vids := mapset.NewSet[uint]()
	tagCondition := c.Param("tags")
	textCondition := c.Param("text")
	sTags := strings.Split(tagCondition, ",") //将用","链接的tags字符串转化为切片
	nTags := uint(len(sTags))                 //统计tag种类的总数
	for _, i := range sTags {
		db.In("text", i).Get(&tagModel)
		db.In("tid", tagModel.Id).Find(&tags)
		for _, j := range tags {
			vids.Add(uint(j.Pid))
			vidsCount[uint(j.Pid)]++ //如果视频在tags出现一次就+1

		}
	}
	//bucket sort
	for i := range vids.Iter() {
		if vidsCount[i] == nTags { //如果视频出现在每一种tag中，vidsCount对应项与tag种类总数相等
			db.ID(i).Get(&video)
			if strings.Contains(video.Title, textCondition) || strings.Contains(video.Description, textCondition) || tagCondition == "" {
				//判断是否又对应text部分,如果没有text部分，就忽略text部分
				videoReturn.Id = video.Id
				videoReturn.Title = video.Title
				videoReturn.CoverPath = video.CoverPath
				videoReturn.Views = video.Views
				videosReturn = append(videosReturn, videoReturn)
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"video": videosReturn,
	})
}
