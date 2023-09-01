package main

import (
	"net/http"
	"strings"

	//"fmt"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-gonic/gin"
)

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

func SearchFourms(c *gin.Context) {
	var fourm Forum
	var forumsC []ForumComment
	var fourmsTitle []Forum //这里是检索标题后的结果
	var forumSearch SearchFourmReturn
	var forumsS []SearchFourmReturn //返回的模板
	ids := mapset.NewSet[uint]()
	textCondition := c.Param("text")
	db.Where("Text like ?", "%"+textCondition+"%").Find(&forumsC)
	for _, i := range forumsC {
		if !ids.Contains(uint(i.Mid)) { //排除同一主帖子下子帖子反复出现关键词
			db.ID(i.Mid).Get(&fourm)
			forumSearch.Id = fourm.Id
			forumSearch.Text = i.Text
			forumSearch.Title = fourm.Title
			forumSearch.Kind = fourm.Kind
			ids.Add(uint(i.Mid))
			forumsS = append(forumsS, forumSearch)
		}
	}
	db.Where("Title like ?", "%"+textCondition+"%").Find(&fourmsTitle)
	for _, i := range fourmsTitle {
		if !ids.Contains(uint(i.Id)) {
			forumSearch.Id = i.Id
			forumSearch.Title = i.Title
			forumSearch.Kind = i.Kind
			forumsS = append(forumsS, forumSearch)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"forum": forumsS,
	})

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
			if j.Kind == 1 {
				vids.Add(uint(j.Pid))
				vidsCount[uint(j.Pid)]++ //如果视频在tags出现一次就+1
			}
		}
	}
	for i := range vids.Iter() {
		if vidsCount[i] == nTags { //如果视频出现在每一种tag中，vidsCount对应项与tag种类总数相等
			db.ID(i).Get(&video)
			if strings.Contains(video.Title, textCondition) || strings.Contains(video.Description, textCondition) {
				//判断是否又对应text部分
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
