package main

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ContainCCA(kinds []uint8, postKind uint8) bool {
	for _, i := range kinds {
		if i == postKind {
			return true
		}
	}
	return false
}

func CheckCircleAvaible(c *gin.Context) {
	strKind := c.Param("type")
	kind, err := strconv.Atoi(strKind)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	if kind > 3 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	session := sessions.Default(c)
	author := session.Get("userid")
	uauthor := int(author.(int64))
	var authorOfCircle []MemberOfCircle
	var circles []gin.H
	has, _ := db.In("uid", uauthor).Count(&MemberOfCircle{})
	if has == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.In("uid", uauthor).Asc("permission").Find(&authorOfCircle)
	var circle Circle
	for _, i := range authorOfCircle {
		db.ID(i.Cid).Get(&circle)
		if i.Permission <= 2 && ContainCCA(circle.Kinds, uint8(kind)) { //need to add check creating type
			circles = append(circles, gin.H{
				"name": circle.Name,
				"id":   i.Cid,
			})
		}
	}
	if len(circles) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSONP(http.StatusOK, gin.H{
			"circles": circles,
		})
	}
}
