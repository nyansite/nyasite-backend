package main

import (
	"net/http"
	"time"

	//"strconv"

	"github.com/gin-gonic/gin"
)

func GetVideoSubscribe(c *gin.Context) {
	userid := GetUserIdWithoutCheck(c)
	var user User
	var subscribes []MemberOfCircle
	var newVideosDisplay []VideoReturn
	db.In("uid", userid).Find(&subscribes)
	db.ID(userid).Get(&user)
	for _, i := range subscribes {
		var videosofSingleCircle []Video
		db.In("author", i.Cid).Find(&videosofSingleCircle)
		for _, j := range videosofSingleCircle {
			if j.CreatedAt >= user.LTCM-302400 {
				var newVideoDisplay VideoReturn
				author := DBGetCircleDataShow(j.Author)
				newVideoDisplay.Author.Id = author.Id
				newVideoDisplay.Author.Name = author.Name
				newVideoDisplay.CoverPath = j.CoverPath
				newVideoDisplay.CreatedAt = j.CreatedAt
				newVideoDisplay.Id = j.Id
				newVideoDisplay.Likes = j.Likes
				newVideoDisplay.Title = j.Title
				newVideoDisplay.Views = j.Views
				newVideosDisplay = append(newVideosDisplay, newVideoDisplay)
			}
		}
	}
	SortVideo(newVideosDisplay, func(p, q *VideoReturn) bool {
		return p.Id > q.Id
	})
	c.JSON(http.StatusOK, gin.H{
		"videos": newVideosDisplay,
	})
	user.LTCM = int(time.Now().Unix())
	db.ID(userid).Cols("l_t_c_m").Update(&user)
	return
}
