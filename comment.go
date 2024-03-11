package main

import (
	//"log"
	"math"
	"net/http"
	"strconv"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-gonic/gin"
)

func BrowseVideoComments(ctx *gin.Context) {
	author := GetUserIdWithCheck(ctx)
	vvid := ctx.Param("id")
	vid, _ := strconv.Atoi(vvid)
	vpg := ctx.Param("pg")
	pg, _ := strconv.Atoi(vpg)
	if pg < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	count, err := db.In("vid", vid).Count(&VideoComment{})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err) //500,正常情况下不会出现
		return
	}
	if count < 1 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if pg > int(math.Ceil(float64(count)/20)) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	userIds := mapset.NewSet[int]()
	comments := DBgetVideoComments(vid, pg, author)
	var userDataShows []UserDataShow
	for _, i := range comments {
		if !userIds.Contains(i.Author) {
			userIds.Add(i.Author)
			userDataShows = append(userDataShows, DBGetUserDataShow(i.Author)) //from user.go
			for _, j := range i.CRdisplay {
				if !userIds.Contains(j.Author) {
					userIds.Add(i.Author)
					userDataShows = append(userDataShows, DBGetUserDataShow(j.Author))
				}
			}
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Body":     comments,
		"UserShow": userDataShows,
		"Count":    count, //总数
	})

}
func BrowseVideoCommentReplies(ctx *gin.Context) {
	author := GetUserIdWithCheck(ctx)
	vcid := ctx.Param("id")
	cid, _ := strconv.Atoi(vcid)
	var comment VideoComment
	has, err := db.ID(cid).Get(&comment)
	var emojiRecord VideoCommentEmojiRecord
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //400
		return
	}
	if !has{
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	has1, _ := db.Where("author = ? AND cid = ?", author, comment.Id).Get(&emojiRecord)
	if !has1 {
		comment.Choose = 0
	} else {
		comment.Choose = emojiRecord.Emoji
	}
	userIds := mapset.NewSet[int]()
	commentrReplies := DBgetVideoCommentReplies(cid, author)
	var userDataShows []UserDataShow
	userIds.Add(comment.Author)
	userDataShows = append(userDataShows, DBGetUserDataShow(comment.Author))
	for _, i := range commentrReplies {
		if !userIds.Contains(i.Author) {
			userIds.Add(i.Author)
			userDataShows = append(userDataShows, DBGetUserDataShow(i.Author))
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Origin":   comment,
		"Body":     commentrReplies,
		"UserShow": userDataShows,
	})
}
func DBgetVideoComments(vid int, page int, author int) []VideoComment {
	var comments []VideoComment
	var commentsReturn []VideoComment
	db.In("vid", vid).Desc("id").Limit(20, (page-1)*20).Find(&comments)
	for _, i := range comments {
		var emojiRecord VideoCommentEmojiRecord
		has, _ := db.Where("author = ? AND cid = ?", author, i.Id).Get(&emojiRecord)
		if !has {
			i.Choose = 0
		} else {
			i.Choose = emojiRecord.Emoji
		}
		//从1开始计数，所以默认-1
		i.CRdisplay = DBgetVideoCommentRepliesShow(int(i.Id), author)
		commentsReturn = append(commentsReturn, i)
	}
	return commentsReturn
}

func DBgetVideoCommentReplies(cid int, author int) []VideoCommentReply {
	var commentReplies []VideoCommentReply
	var commentRepliesReturn []VideoCommentReply
	db.In("cid", cid).Desc("id").Find(&commentReplies)
	for _, i := range commentReplies {
		exist, _ := db.Where("author = ? and crid = ?", author, i.Id).Count(&VideoCommentReplyLikeRecord{})
		if exist != 0 {
			i.Like_c = true
		} else {
			i.Like_c = false
		}
		i.Likes--
		commentRepliesReturn = append(commentRepliesReturn, i)
	}
	return commentRepliesReturn
}
func DBgetVideoCommentRepliesCount(cid int) int {
	count, _ := db.In("cid", cid).Count(&VideoCommentReply{})
	return int(math.Ceil(float64(count) / 20))
}

func DBgetVideoCommentRepliesShow(cid int, author int) []VideoCommentReply {
	var commentReplies []VideoCommentReply
	var commentRepliesReturn []VideoCommentReply
	db.In("cid", cid).Desc("id").Limit(3, 0).Find(&commentReplies)
	for _, i := range commentReplies {
		exist, _ := db.Where("author = ? AND crid = ?", author, i.Id).Count(&VideoCommentReplyLikeRecord{})
		if exist != 0 {
			i.Like_c = true
		} else {
			i.Like_c = false
		}
		i.Likes--
		commentRepliesReturn = append(commentRepliesReturn, i)
	}
	return commentRepliesReturn
}

func AddVideoComment(ctx *gin.Context) {
	uauthor := GetUserIdWithoutCheck(ctx)
	vid, text := ctx.PostForm("vid"), ctx.PostForm("text")
	vvid, _ := strconv.Atoi(vid)
	uvid := int(vvid)
	cid := DBaddVideoComment(uvid, uauthor, text)
	ctx.String(http.StatusOK, "%v", cid)
}

func DBaddVideoComment(vid int, author int, text string) int {
	comment := VideoComment{Vid: vid, Text: text, Author: author,
		Like: 1, Dislike: 1, Smile: 1, Celebration: 1, Confused: 1, Heart: 1, Rocket: 1, Eyes: 1}
	db.Insert(&comment)
	return int(comment.Id)
}

func AddVideoCommentReply(ctx *gin.Context) {
	uauthor := GetUserIdWithoutCheck(ctx)
	cid, text := ctx.PostForm("cid"), ctx.PostForm("text")
	vcid, _ := strconv.Atoi(cid)
	ucid := int(vcid)
	crid := DBaddVideoCommentReply(ucid, uauthor, text)
	ctx.String(http.StatusOK, "%v", crid)
}

func DBaddVideoCommentReply(cid int, author int, text string) int {
	commentReply := VideoCommentReply{Cid: cid, Text: text, Author: author, Likes: 1}
	db.Insert(&commentReply)
	return int(commentReply.Id)
}

func ClikckCommentEmoji(ctx *gin.Context) {
	uauthor := GetUserIdWithoutCheck(ctx)
	emoji, cid := ctx.PostForm("emoji"), ctx.PostForm("cid")
	vemoji, _ := strconv.Atoi(emoji)
	uemoji := int8(vemoji)
	vcid, _ := strconv.Atoi(cid)
	exist, _ := db.Where("author = ? and cid = ?", uauthor, cid).Count(&VideoCommentEmojiRecord{})
	if vemoji > 8 || vemoji < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest) //传入的表情编号>7(不存在)
		return
	}
	if exist == 0 {
		DBaddCommentEmoji(vcid, uemoji, uauthor)
		return
	} else {
		var existedEmojiRecord VideoCommentEmojiRecord
		db.Where("author = ? and cid = ?", uauthor, cid).Get(&existedEmojiRecord)
		if existedEmojiRecord.Emoji == uemoji {
			DBdeleteCommentEmoji(vcid, uemoji, uauthor)
			return
		} else {
			DBchangeCommentEmoji(vcid, uemoji, existedEmojiRecord)
			return
		}
	}
}

func DBaddCommentEmoji(cid int, emoji int8, author int) {
	var comment VideoComment
	db.ID(cid).Get(&comment)
	switch emoji {
	case 1:
		comment.Like++
	case 2:
		comment.Dislike++
	case 3:
		comment.Smile++
	case 4:
		comment.Celebration++
	case 5:
		comment.Confused++
	case 6:
		comment.Heart++
	case 7:
		comment.Rocket++
	case 8:
		comment.Eyes++
	}
	emojiRecord := VideoCommentEmojiRecord{Author: author, Cid: cid, Emoji: emoji}
	db.Insert(&emojiRecord)
	db.ID(cid).Update(&comment)
}

func DBchangeCommentEmoji(cid int, emoji int8, emojiRecord VideoCommentEmojiRecord) {
	//直接导入查询到的表情记录，避免重复查询
	var comment VideoComment
	db.ID(cid).Get(&comment)
	oEmoji := emojiRecord.Emoji
	switch oEmoji {
	case 1:
		comment.Like--
	case 2:
		comment.Dislike--
	case 3:
		comment.Smile--
	case 4:
		comment.Celebration--
	case 5:
		comment.Confused--
	case 6:
		comment.Heart--
	case 7:
		comment.Rocket--
	case 8:
		comment.Eyes--
	}
	switch emoji {
	case 1:
		comment.Like++
	case 2:
		comment.Dislike++
	case 3:
		comment.Smile++
	case 4:
		comment.Celebration++
	case 5:
		comment.Confused++
	case 6:
		comment.Heart++
	case 7:
		comment.Rocket++
	case 8:
		comment.Eyes++
	}
	emojiRecord.Emoji = emoji
	db.ID(cid).Update(&comment)
	db.Where("author = ? and cid = ?", emojiRecord.Author, cid).Cols("emoji").Update(&emojiRecord)
}

func DBdeleteCommentEmoji(cid int, emoji int8, author int) {
	var comment VideoComment
	db.ID(cid).Get(&comment)
	switch emoji {
	case 1:
		comment.Like--
	case 2:
		comment.Dislike--
	case 3:
		comment.Smile--
	case 4:
		comment.Celebration--
	case 5:
		comment.Confused--
	case 6:
		comment.Heart--
	case 7:
		comment.Rocket--
	}
	db.ID(cid).Update(&comment)
	db.Where("author = ? and cid = ?", author, cid).Delete(&VideoCommentEmojiRecord{})
}

func ClickCommentReplyLike(ctx *gin.Context) {
	uauthor := GetUserIdWithoutCheck(ctx)
	crid := ctx.PostForm("crid")
	vcrid, _ := strconv.Atoi(crid)
	count, _ := db.Where("author = ? and crid = ?", uauthor, crid).Count(&VideoCommentReplyLikeRecord{})
	if count == 0 {
		DBaddCommentReplyLike(uauthor, vcrid)
	} else {
		DBdeleteCommentReplyLike(uauthor, vcrid)
	}
}
func DBaddCommentReplyLike(author int, crid int) {
	var commentReply VideoCommentReply
	db.ID(crid).Get(&commentReply)
	println(commentReply.Likes)
	commentReply.Likes++
	commentReplyLikeRecord := VideoCommentReplyLikeRecord{Author: author, Crid: crid}
	db.Insert(&commentReplyLikeRecord)
	db.ID(crid).Update(&commentReply)
}

func DBdeleteCommentReplyLike(author int, crid int) {
	var commentReply VideoCommentReply
	db.ID(crid).Get(&commentReply)
	println(commentReply.Likes)
	commentReply.Likes--
	db.Where("author = ? and crid = ?", author, crid).Delete(&VideoCommentReplyLikeRecord{})
	db.ID(crid).Update(&commentReply)
}
