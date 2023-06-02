package main

import (
	"gorm.io/gorm"
)

func DBaddMainForum(text string, title string, author uint, kind uint8, db *gorm.DB) {
	var mainForum MainForum
	mainForum.Title = title
	mainForum.Author = author
	mainForum.Views = 0
	db.Create(&mainForum)
	println(text)
	println(kind)
}
