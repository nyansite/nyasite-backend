package main

func DBaddMainForum(text string, isMD bool, title string, author uint) uint{
	mainForum := &Forum{Title: title, Author: author}
	mainForum.Comment = append(mainForum.Comment, ForumComment{Text: text, IsMD: isMD, Author: author})
	db.Create(&mainForum)
	return mainForum.ID
}

func DBaddComment(text string, isMD bool, uid uint, author uint) uint{
	var forum Forum
	db.Take(&forum, uid)
	forum.Comment = append(forum.Comment, ForumComment{Text: text, IsMD: isMD, Author: author})
	db.Save(&forum)
	return forum.ID
}
