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

func addComment(str string, uid uint, cid uint, author uint) { //测试用
	var utilForum UtilForum
	db.Preload("CommentP").First(&utilForum, uid)

	// fmt.Println(video.CommentP[len(video.CommentP)-1].ID)
	var com []Comment
	db.Find(&com, "Pid = ?", utilForum.CommentP[len(utilForum.CommentP)-1].ID)

	if len(com) >= 16 {
		utilForum.CommentP = append(utilForum.CommentP, CommentPage{uid: utilForum.ID})

	}
	pg := utilForum.CommentP[len(utilForum.CommentP)-1]

	pg.Comment = append(pg.Comment, Comment{Text: str, Cid: cid, Author: author, Pid: mid})
	db.Save(&pg)
}
