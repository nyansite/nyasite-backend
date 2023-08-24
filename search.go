package main

import (
	"net/http"
	//"strings"

	"fmt"

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
	var forumSearch SearchFourmsReturn
	var forumsS []SearchFourmsReturn
	ids := mapset.NewSet[uint]()
	textCondition := c.Param("text")
	db.Where("Text like ?", "%"+textCondition+"%").Find(&forumsC)
	for _, i := range forumsC {
		fmt.Println(ids)
		fmt.Println(ids.Contains(i.Mid))
		if !ids.Contains(i.Mid) {
			db.ID(i.Mid).Get(&fourm)
			forumSearch.Id = fourm.Id
			forumSearch.Text = i.Text
			forumSearch.Title = fourm.Title
			forumSearch.Kind = fourm.Kind
			fmt.Println(i.Mid)
			ids.Add(i.Mid)
			forumsS = append(forumsS, forumSearch)
		}
	}
	db.Where("Title like ?", "%"+textCondition+"%").Find(&fourmsTitle)
	for _, i := range fourmsTitle {
		fmt.Println(ids.Contains(uint(i.Id)))
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

/*
func SearchVideos(c *gin.Context) {
	var tagModel TagModel
	var tags []Tag
	vids := mapset.NewSet[uint]()
	tagCondition := c.Param("tags")
	textCondition := c.Param("text")
	sTags := strings.Split(tagCondition, ",")
	for _, i := range sTags {
		db.In("text", i).Find(&tagModel)
		db.In("tid", tagModel.Id).Find(&tags)
		for _, j := range tags {
			vids.Add(uint(j.Id))
		}
	}
	for _, i := range vids.Iter() {

	}
}
*/
