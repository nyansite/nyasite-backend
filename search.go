package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchTag(c *gin.Context) {
	tagConditon := c.Param("tc")
	var tagMs []TagModel
	var tagList []string
	db.Where("Text like ?", tagConditon+"%").Find(&tagMs)
	for _, i := range tagMs {
		tagList = append(tagList, i.Text)
	}
	c.JSON(http.StatusOK, gin.H{"taglist": tagList})
	return
}
