package main

func DBaddMainForum(text string, title string, author uint, ismd bool) {
	mainForum := Forum{Title: title, Author: author}
	mainForum.Comment = append(mainForum.Comment, ForumComment{Text: text, Author: author, IsMD: ismd})
	db.Create(&mainForum)
	return
}

func DBaddUtilForum(text string, fid uint, author uint, ismd bool) {
	var mainForum Forum
	db.Take(&mainForum, fid)
	mainForum.Comment = append(mainForum.Comment, ForumComment{Text: text, IsMD: ismd, Author: author})
	db.Save(&mainForum)
	return
}
