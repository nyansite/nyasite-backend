package main

func DBaddMainForum(text string, title string, author uint, ismd bool) {
	forum := Forum{Title: title, Author: author}
	forum.Comment = append(forum.Comment, ForumComment{Text: text, Author: author, IsMD: ismd})
	db.Insert(&forum)
	return
}

func DBaddUtilForum(text string, fid uint, author uint, ismd bool) {
	var forum Forum
	db.ID(fid).Get(&forum)
	forum.Comment = append(forum.Comment, ForumComment{Text: text, IsMD: ismd, Author: author})
	db.Update(&forum)
	return
}
