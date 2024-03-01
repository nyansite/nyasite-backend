package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CheckPremissionOfCircle(c *gin.Context) { //for manage page
	cidStr := c.Param("cid")
	cid, _ := strconv.Atoi(cidStr)
	uid := GetUserIdWithoutCheck(c)
	var memberOfCircle MemberOfCircle
	has, _ := db.Where("cid = ? and uid = ?", cid, uid).Count(&MemberOfCircle{})
	if has == 0 {
		c.AbortWithStatus(http.StatusForbidden)
	}
	db.Where("cid = ? and uid = ?", cid, uid).Get(&memberOfCircle)
	if memberOfCircle.Permission > 0 {
		c.String(http.StatusOK, "%v", memberOfCircle.Permission)
	} else {
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func GetAllMembersOfCircle(c *gin.Context) {
	userid := GetUserIdWithoutCheck(c)
	cidStr := c.Param("cid")
	cid, _ := strconv.Atoi(cidStr)
	var members []MemberOfCircle
	var membersDisplay []UserDataShow
	permission := DBgetRelationToCircle(cid, c)
	if permission < 1 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	db.Where("permission > 0 and cid = ? and uid <> ?", cid, userid).Find(&members)
	for _, i := range members {
		memberDisplay := DBGetUserDataShow(i.Uid)
		membersDisplay = append(membersDisplay, memberDisplay)
	}
	self := DBGetUserDataShow(userid)
	c.JSON(http.StatusOK, gin.H{
		"self":   self,
		"others": membersDisplay,
	})
}

func InviteMember(c *gin.Context) {
	inviteeIdStr := c.PostForm("eid")
	inviterIdStr := c.PostForm("rid")
	circleIdStr := c.PostForm("cid")
	inviteeId, _ := strconv.Atoi(inviteeIdStr)
	inviterId, _ := strconv.Atoi(inviterIdStr)
	circleId, _ := strconv.Atoi(circleIdStr)
	invitation := Invitation{
		Inviter: inviterId,
		Invitee: inviteeId,
		Circle:  circleId,
		stauts:  false,
	}
	db.InsertOne(&invitation)
}
