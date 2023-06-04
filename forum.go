package main

func DBaddMainForum(text string, title string, author uint, ismd bool) {
	mainForum := MainForum{Title: title, Views: 0, Author: author}
	db.Create(&mainForum)
	mainForum.UnitP = append(mainForum.UnitP, UnitForum{Text: text, Author: author,
		Mid: mainForum.ID, Cid: 0, IsMD: ismd, Likes: 0})
	db.Save(&mainForum)
	return
}

func DBaddUtilForum(text string, mid uint, cid uint, author uint, ismd bool) {
	var mainForum MainForum
	db.Last(&mainForum, mid)
	mainForum.UnitP = append(mainForum.UnitP, UnitForum{Text: text, Mid: mid,
		Cid: cid, IsMD: ismd, Author: author, Likes: 0, Dislikes: 0})
	db.Save(&mainForum)
	return
}

func DBaddComment(text string, uid uint, author uint) {
	var unitForum UnitForum
	db.Last(&unitForum, uid)
	unitForum.CommentP = append(unitForum.CommentP, Comment{Text: text, Uid: uid, Author: author})
	db.Save(&unitForum)
	return
}

func DBlike(id uint) {
	var unitForum UnitForum
	db.Last(&unitForum, id)
	unitForum.Likes = unitForum.Likes + 1
	db.Save(&unitForum)
	return
}

func DBdislike(id uint) {
	var unitForum UnitForum
	db.Last(&unitForum, id)
	unitForum.Dislikes = unitForum.Dislikes + 1
	db.Save(&unitForum)
	return
}
