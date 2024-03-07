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
