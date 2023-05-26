package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func addMainForum(text string, title string, author uint) {
	mainForum := MainForum{Title: title, Author: author, Views: 0}
	db.Create(&mainForum)
	mainForum.UnitP = append(mainForum.UnitP, UtilForumPage{Tid: mainForum.ID})
	mainForum.UnitP[0].UtilForum = append(mainForum.UnitP[0].UtilForum, UtilForum{Pid: 0, Text: text, Cid: 0, Author: author})
	db.Save(&mainForum)
}
func addUtilForum(text string, mid uint, cid uint, author uint) {
	var mainForum MainForum
	db.Preload("UnitP").First(&mainForum, mid)
	var com []UtilForum
	db.Find(&com, "Pid = ?", mainForum.UnitP[len(mainForum.UnitP)-1].ID)

	if len(com) >= 16 {
		mainForum.UnitP = append(mainForum.UnitP, UtilForumPage{Pid: mainForum.ID})

	}
	pg := mainForum.UnitP[len(mainForum.UnitP)-1]

	pg.UtilForum = append(pg.UtilForum, UtilForum{Text: text, Cid: cid, Author: author, Pid: mid})
	db.Save(&pg)
}
