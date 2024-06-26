package main

import (
	"net/http"
	"sort"
	"time"

	"strconv"

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
}

func GetCircleAffairs(c *gin.Context) {
	userid := GetUserIdWithoutCheck(c)
	timeLimit := time.Now().Unix() - 31556952
	var invitationsSelf []Invitation
	var invitationsCircle []Invitation
	var messages []CircleAffairMessage
	circlesManagedId := DBgetCirclesRelatedTo(userid)
	for _, i := range circlesManagedId {

		var invitationsCircleUnit []Invitation
		db.In("circle", i).Where("created_at > ?", timeLimit).Find(&invitationsCircleUnit)
		invitationsCircle = append(invitationsCircle, invitationsCircleUnit...)
		var membersOfCircle []MemberOfCircle
		db.Where("cid = ? and uid <> ? and updated_at > ?", i, userid, timeLimit).Find(&membersOfCircle)
		for _, j2 := range membersOfCircle {
			var message CircleAffairMessage
			newcomer := DBGetUserDataShow(j2.Uid)
			message.ReciverId = j2.Uid
			message.ReciverName = newcomer.Name
			circle := DBGetCircleDataShow(j2.Cid)
			message.CircleName = circle.Name
			message.CircleId = j2.Cid
			message.Kind = 0
			messages = append(messages, message)
		}
	}
	for _, i := range invitationsCircle {
		if !i.Stauts {
			var message CircleAffairMessage
			inviter := DBGetUserDataShow(i.Inviter)
			message.SenderId = i.Inviter
			message.SenderName = inviter.Name
			invitee := DBGetUserDataShow(i.Invitee)
			message.ReciverId = i.Invitee
			message.ReciverName = invitee.Name
			circle := DBGetCircleDataShow(i.Circle)
			message.CircleId = i.Circle
			message.CircleName = circle.Name
			message.Time = i.CreatedAt
			message.Kind = i.Kind + 5
			messages = append(messages, message)
		} else {
			var message CircleAffairMessage
			inviter := DBGetUserDataShow(i.Inviter)
			message.SenderId = i.Inviter
			message.SenderName = inviter.Name
			invitee := DBGetUserDataShow(i.Invitee)
			message.ReciverId = i.Invitee
			message.ReciverName = invitee.Name
			circle := DBGetCircleDataShow(i.Circle)
			message.CircleId = i.Circle
			message.CircleName = circle.Name
			message.Time = i.CreatedAt
			message.Kind = 9
			messages = append(messages, message)
		}
	}
	db.Where("invitee = ? and created_at > ? and stauts = false", userid, timeLimit).Find(&invitationsSelf)
	for _, i := range invitationsSelf {
		var message CircleAffairMessage
		inviter := DBGetUserDataShow(i.Inviter)
		message.SenderId = i.Inviter
		message.SenderName = inviter.Name
		invitee := DBGetUserDataShow(i.Invitee)
		message.ReciverId = i.Invitee
		message.ReciverName = invitee.Name
		circle := DBGetCircleDataShow(i.Circle)
		message.CircleId = i.Circle
		message.CircleName = circle.Name
		message.Time = i.CreatedAt
		message.Kind = i.Kind + 2
		messages = append(messages, message)
	}
	sort.Sort(CircleAffairsSliceDecrement(messages))
	for i1 := range messages {
		messages[i1].Id = i1
	}
	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func ReplyInvitation(c *gin.Context) {
	inviteeIdStr := c.PostForm("eid")
	circleIdStr := c.PostForm("cid")
	stautsStr := c.PostForm("stauts")
	inviteeId, _ := strconv.Atoi(inviteeIdStr)
	circleId, _ := strconv.Atoi(circleIdStr)
	stauts, _ := strconv.ParseBool(stautsStr)
	var invitation Invitation

	hasInvitation, _ := db.In("invitee", inviteeId).In("circle", circleId).In("stauts", false).Get(&invitation)
	if !hasInvitation{
		c.AbortWithStatus(http.StatusFailedDependency)
		return
	}
	if stauts {
		var memberOfCircle MemberOfCircle
		hasMemberOfCircle, _ := db.In("uid", inviteeId).In("cid", circleId).Get(&memberOfCircle)
		if !hasMemberOfCircle{
			memberOfCircle.Uid = inviteeId
			memberOfCircle.Cid = circleId
			memberOfCircle.Permission = invitation.Kind
			db.InsertOne(&memberOfCircle)
		} else {
			memberOfCircle.Permission = invitation.Kind
			db.In("uid", inviteeId).In("cid", circleId).Cols("permission").Update(&memberOfCircle)
		}
		db.In("invitee", inviteeId).In("circle", circleId).Delete(&Invitation{})
	} else {
		invitation.Stauts = true
		db.In("invitee", inviteeId).In("circle", circleId).Cols("stauts").Update(&invitation)
	}
}

func GetCheckMessage(c *gin.Context) {
	userid := GetUserIdWithoutCheck(c)
	circlesId := DBgetCirclesRelatedTo(userid)
	timeLimit := time.Now().Unix() - 31556952
	var messages []CheckMessage
	for _, i := range circlesId {
		var videosPast []Video
		db.Where("author = ? and updated_at > ?", i, timeLimit).Find(&videosPast)
		for _, j1 := range videosPast {
			var message CheckMessage
			message.Image = j1.CoverPath
			message.Name = j1.Title
			message.Kind = 0
			message.Time = j1.CreatedAt
			messages = append(messages, message)
		}
		var videosRejected []VideoNeedToCheck
		db.Where("author = ? and stauts = true and updated_at > ?", i, timeLimit).Find(&videosRejected)
		for _, j2 := range videosRejected {
			var message CheckMessage
			message.Image = j2.CoverPath
			message.Name = j2.Title
			message.Kind = 1
			message.Time = j2.UpdatedAt
			message.Reason = j2.Reason
			message.DBId = int(j2.Id)
			messages = append(messages, message)
		}

	}
	var asOwnerOfCircle []MemberOfCircle
	db.Where("uid = ? and permission = 4 and updated_at > ?", userid, timeLimit).Find(&asOwnerOfCircle)
	for _, j3 := range asOwnerOfCircle {
		var circlePast Circle
		db.ID(j3.Cid).Get(&circlePast)
		var message CheckMessage
		message.Image = circlePast.Avatar
		message.Name = circlePast.Name
		message.Kind = 2
		message.Time = circlePast.CreatedAt
		messages = append(messages, message)
	}
	var CircleRejected []ApplyCircle
	db.Where("applicant = ? and stauts = true and updated_at > ?", userid, timeLimit).Find(&CircleRejected)
	for _, j4 := range CircleRejected {
		var message CheckMessage
		message.Image = j4.Avatar
		message.Name = j4.Name
		message.Kind = 3
		message.Time = j4.UpdatedAt
		messages = append(messages, message)
	}
	sort.Sort(CheckMessageSliceDecrement(messages))
	for i1 := range messages {
		messages[i1].Id = i1
	}
	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func DeleteVideoNeedToCheck(c *gin.Context) {
	vcid := c.PostForm("id")
	db.ID(vcid).Delete(&VideoNeedToCheck{})
}
