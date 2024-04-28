package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SearchUsers(c *gin.Context) {
	userid := GetUserIdWithoutCheck(c)
	clipOfName := c.Param("name")
	var users []User
	var usersDisplay []gin.H
	db.Where("name LIKE ?", "%"+clipOfName+"%").Find(&users)
	for _, i := range users {
		if int(i.Id) != userid {
			usersDisplay = append(usersDisplay, gin.H{
				"id":     i.Id,
				"avatar": i.Avatar,
				"name":   i.Name,
			})
		}
	}
	if len(usersDisplay) > 5 {
		usersDisplay = usersDisplay[:5]
	}
	c.JSON(http.StatusOK, gin.H{
		"users": usersDisplay,
	})
}

func CheckPremissionOfCircle(c *gin.Context) { //for manage page
	cidStr := c.Param("cid")
	cid, _ := strconv.Atoi(cidStr)
	uid := GetUserIdWithoutCheck(c)
	var memberOfCircle MemberOfCircle
	has, _ := db.In("cid", cid).In("uid", uid).Count(&MemberOfCircle{})
	if has == 0 {
		c.AbortWithStatus(http.StatusForbidden)
	}
	db.In("cid", cid).In("uid", uid).Get(&memberOfCircle)
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
	var membersDisplay []UserDataShowWithPermission
	permission := DBgetRelationToCircle(cid, c)
	if permission < 1 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	db.Where("permission > 0 and cid = ? and uid <> ?", cid, userid).Find(&members)
	for _, i := range members {
		var memberDisplay UserDataShowWithPermission
		memberData := DBGetUserDataShow(i.Uid)
		memberDisplay.Avatar = memberData.Avatar
		memberDisplay.Name = memberData.Name
		memberDisplay.Id = memberData.Id
		memberDisplay.Permission = i.Permission
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
	circleIdStr := c.PostForm("cid")
	kindStr := c.PostForm("kind")
	inviteeId, _ := strconv.Atoi(inviteeIdStr)
	inviterId := GetUserIdWithoutCheck(c)
	circleId, _ := strconv.Atoi(circleIdStr)
	kind, _ := strconv.Atoi(kindStr)
	privilege := DBgetRelationToCircle(circleId, c)
	switch privilege {
	case 3:
		if (kind == 1) || (kind == 2) {
			invitation := Invitation{
				Inviter: inviterId,
				Invitee: inviteeId,
				Circle:  circleId,
				Kind:    uint8(kind),
				Stauts:  false,
			}
			db.InsertOne(&invitation)
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		return
	case 4:
		if (kind >= 1) && (kind <= 3) {
			invitation := Invitation{
				Inviter: inviterId,
				Invitee: inviteeId,
				Circle:  circleId,
				Kind:    uint8(kind),
				Stauts:  false,
			}
			db.InsertOne(&invitation)
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		return
	default:
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}

func KickOut(c *gin.Context) {
	selfId := GetUserIdWithoutCheck(c)
	strUid := c.PostForm("uid")
	strCid := c.PostForm("cid")
	cid, _ := strconv.Atoi(strCid)
	uid, _ := strconv.Atoi(strUid)
	var memberKickedOut MemberOfCircle
	var memberSelf MemberOfCircle
	db.In("uid", uid).In("cid", cid).Get(&memberKickedOut)
	db.In("uid", selfId).In("cid", cid).Get(&memberSelf)
	if (uid == 0) && (memberSelf.Permission != 4) {
		db.In("uid", selfId).In("cid", cid).Delete(&MemberOfCircle{})
		kickOutRecord := Discharge{
			WhoDischarge:     selfId,
			WhoBeDischargeed: selfId,
			Circle:           cid,
		}
		db.InsertOne(&kickOutRecord)
		return
	} else if ((memberSelf.Permission == 4) && (memberKickedOut.Permission <= 3) && (uid != 0)) ||
		((memberSelf.Permission == 3) && (memberKickedOut.Permission <= 2)) {
		db.In("uid", uid).In("cid", cid).Delete(&MemberOfCircle{})
		kickOutRecord := Discharge{
			WhoDischarge:     selfId,
			WhoBeDischargeed: uid,
			Circle:           cid,
		}
		db.InsertOne(&kickOutRecord)
		return
	} else {
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func DeleteVideo(c *gin.Context) {
	vidStr := c.PostForm("vid")
	vid, _ := strconv.Atoi(vidStr)
	userid := GetUserIdWithoutCheck(c)
	var video Video
	db.ID(vid).Get(&video)
	if video.Upid == userid {
		db.ID(vid).Delete(&Video{})
		return
	} else {
		var memberOfCircle MemberOfCircle
		db.In("uid", userid).In("cid", video.Author).Get(&memberOfCircle)
		if memberOfCircle.Permission >= 3 {
			db.ID(vid).Delete(&Video{})
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}
