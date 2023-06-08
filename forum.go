package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DBaddMainForum(text string, title string, author uint, ismd bool) {
	forum := Forum{Title: title, Author: author}
	// forum.Comment = append(forum.Comment, ForumComment{Text: text, Author: author, IsMD: ismd})
	db.Insert(&forum)
	return
}

func DBaddUtilForum(text string, fid uint, author uint, ismd bool) {
	var forum Forum
	db.ID(fid).Get(&forum)
	// forum.Comment = append(forum.Comment, ForumComment{Text: text, IsMD: ismd, Author: author})
	db.Update(&forum)
	return
}

func DBaddEmoji(fid uint, id uint, emoji uint) {
	var forum Forum
	db.ID(fid).Get(&forum)
	// forum.Comment[id-1].Emoji[emoji] = forum.Comment[id-1].Emoji[emoji] + 1
	db.Update(&forum)
	return
}

func FindMainForum(fid uint, c *gin.Context) {
	var forum Forum
	db.ID(fid).Get(&forum)
	var user User
	db.ID(forum.Author).Get(&user)
	authorName := user.Name
	c.JSON(http.StatusOK, gin.H{"title": forum.Title,
		"authorID": forum.Author, "author": authorName, "views": forum.Views, "creatTime": forum.CreatedAt})
	return
}
func FindUnitForum(fid uint, id uint, c *gin.Context) {
	var forum Forum
	db.ID(fid).Get(&forum)
	// unitForum := forum.Comment[id-1]
	// var user User
	// db.ID(unitForum.Author).Get(&user)
	// authorName := user.Name
	// c.JSON(http.StatusOK, gin.H{"text": unitForum.Text, "idmd": unitForum.IsMD,
	// 	"authorID": unitForum.Author, "author": authorName, "creatTime": forum.CreatedAt,
	// 	"emoji": unitForum.Emoji})
	return
}
