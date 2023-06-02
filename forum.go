package main

func DBaddMainForum(text string, title string, author uint, kind uint8) {
	mainForum := MainForum{Title: title, Views: 0, Author: author}
	db.Create(&mainForum)
	mainForum.UnitP = append(mainForum.UnitP, UnitForum{Text: text, Author: author,
		Mid: mainForum.ID, Cid: 0, Type: kind, Likes: 0})
	db.Save(&mainForum)
	return
}

func DBaddUtilForum(text string, mid uint, cid uint, kind uint8, author uint) {
	var mainForum MainForum
	db.Last(&mainForum, mid)
	mainForum.UnitP = append(mainForum.UnitP, UnitForum{Text: text, Mid: mid,
		Cid: cid, Type: kind, Author: author, Likes: 0})
	db.Save(&mainForum)
	return
}

func DBaddComment(text string, uid uint, author uint) {
	var unitForum UnitForum
	db.Last(&unitForum, uid)
	unitForum.CommentP = append(unitForum.CommentP, Comment{Text: text, Uid: uid, Likes: 0})
	db.Save(&unitForum)
	return
}
