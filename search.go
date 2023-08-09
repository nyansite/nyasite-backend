package main

import (
	"net/http"

	mapset "github.com/deckarep/golang-set"
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
	titles := mapset.NewSet()
	textCondition := c.Param("tc")
	db.Where("Text like ?", "%"+textCondition+"%").Find(&forumsC)
	for _, i := range forumsC {
		db.ID(i.Mid).Get(&fourm)
		if !titles.Contains(fourm.Title) {
			forumSearch.Id = fourm.Id
			forumSearch.Text = i.Text
			forumSearch.Title = fourm.Title
			forumSearch.Kind = fourm.Kind
			titles.Add(fourm.Title)
			forumsS = append(forumsS, forumSearch)
		}
	}
	db.Where("Title like ?", "%"+textCondition+"%").Find(&fourmsTitle)
	for _, i := range fourmsTitle {
		if !titles.Contains(i.Title) {
			forumSearch.Id = i.Id
			forumSearch.Title = i.Title
			forumSearch.Kind = i.Kind
			forumsS = append(forumsS, forumSearch)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"forum": forumsS,
	})
}
