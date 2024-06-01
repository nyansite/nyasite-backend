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
	db.In("uid", uid).Desc("last_play").Find(&history)
	c.JSON(http.StatusOK, gin.H{
		"Count": count,
		"Body":  history,
	})
}
