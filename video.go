package main

import (
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	UUID "github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func AddVideoTag(c *gin.Context) {
	strVid := c.PostForm("vid")
	vVid, _ := strconv.Atoi(strVid)
	uVid := int(vVid)
	strTagId := c.PostForm("tagid")
	vTagId, _ := strconv.Atoi(strTagId)
	uTagId := int(vTagId)
	DBaddVideoTag(uVid, uTagId)
}

// 视频评论部分

func BrowseVideoComments(ctx *gin.Context) {
	session := sessions.Default(ctx)
	author := session.Get("userid")
	uauthor := int(author.(int64))
	vvid := ctx.Param("id")
	vid, err := strconv.Atoi(vvid)
	vpg := ctx.Param("pg")
	pg, err := strconv.Atoi(vpg)
	var video Video
	if pg < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	pg -= 1
	count, err := db.In("vid", vid).Count(&VideoComment{})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) //500,正常情况下不会出现
		log.Println(err)
		return
	}
	userIds := mapset.NewSet[int]()
	comments := DBgetVideoComments(vid, pg, uauthor)
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
		"Origin":    video,
		"Body":      comments,
		"UserShow":  userDataShows,
		"PageCount": math.Ceil(float64(count) / 20), //总页数
	})

}
func BrowseVideoCommentReplies(ctx *gin.Context) {
	session := sessions.Default(ctx)
	author := session.Get("userid")
	uauthor := int(author.(int64))
	vcid := ctx.Param("id")
	cid, err := strconv.Atoi(vcid)
	vpg := ctx.Param("pg")
	pg, err := strconv.Atoi(vpg)
	var comment VideoComment
	has, err := db.ID(cid).Get(&comment)
	var emojiRecord VideoCommentEmojiRecord
	if pg < 1 || err != nil || has == false {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	count, _ := db.Where("author = ? AND cid = ?", author, comment.Id).Count(&emojiRecord)
	if count == 0 {
		comment.Choose = 0
	} else {
		comment.Choose = emojiRecord.Emoji
	}
	comment.Like--
	comment.Dislike--
	comment.Smile--
	comment.Celebration--
	comment.Confused--
	comment.Heart--
	comment.Rocket--
	comment.Eyes--
	pg -= 1
	count, err1 := db.In("cid", cid).Count(&VideoCommentReply{})
	if err1 != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) //500,正常情况下不会出现
		log.Println(err)
		return
	}
	userIds := mapset.NewSet[int]()
	commentrReplies := DBgetVideoCommentReplies(cid, pg, uauthor)
	var userDataShows []UserDataShow
	for _, i := range commentrReplies {
		if !userIds.Contains(i.Author) {
			userIds.Add(i.Author)
			userDataShows = append(userDataShows, DBGetUserDataShow(i.Author))
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Origin":    comment,
		"Body":      commentrReplies,
		"UserShow":  userDataShows,
		"PageCount": math.Ceil(float64(count) / 20), //总页数
	})
}
func DBgetVideoComments(vid int, page int, author int) []VideoComment {
	var comments []VideoComment
	var commentsReturn []VideoComment
	db.In("vid", vid).Limit(20, (page-1)*20).Find(&comments)
	for _, i := range comments {
		var emojiRecord VideoCommentEmojiRecord
		count, _ := db.Where("author = ? AND cid = ?", author, i.Id).Count(&emojiRecord)
		if count == 0 {
			i.Choose = 0
		} else {
			db.Where("author = ? AND cid = ?", author, i.Id).Get(&emojiRecord)
			i.Choose = emojiRecord.Emoji
		}
		i.Like--
		i.Dislike--
		i.Smile--
		i.Celebration--
		i.Confused--
		i.Heart--
		i.Rocket--
		i.Eyes--
		//从1开始计数，所以默认-1
		i.CRdisplay = DBgetVideoCommentRepliesShow(int(i.Id), author)
		commentsReturn = append(commentsReturn, i)
	}
	return commentsReturn
}

func DBgetVideoCommentReplies(cid int, page int, author int) []VideoCommentReply {
	var commentReplies []VideoCommentReply
	var commentRepliesReturn []VideoCommentReply
	db.In("cid", cid).Limit(20, (page-1)*20).Find(&commentReplies)
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
	db.In("cid", cid).Limit(3, 0).Find(&commentReplies)
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

func AddVideoComment(ctx *gin.Context) {
	session := sessions.Default(ctx)
	author := session.Get("userid")
	uauthor := int(author.(int64))
	vid, text := ctx.PostForm("vid"), ctx.PostForm("text")
	vvid, _ := strconv.Atoi(vid)
	uvid := int(vvid)
	DBaddVideoComment(uvid, uauthor, text)
	return
}

func DBaddVideoComment(vid int, author int, text string) {
	comment := VideoComment{Vid: vid, Text: text, Author: author,
		Like: 1, Dislike: 1, Smile: 1, Celebration: 1, Confused: 1, Heart: 1, Rocket: 1, Eyes: 1}
	db.Insert(comment)
	return
}

func AddVideoCommentReply(ctx *gin.Context) {
	session := sessions.Default(ctx)
	author := session.Get("userid")
	uauthor := int(author.(int64))
	cid, text := ctx.PostForm("cid"), ctx.PostForm("text")
	vcid, _ := strconv.Atoi(cid)
	ucid := int(vcid)
	DBaddVideoCommentReply(ucid, uauthor, text)
	return
}
func DBaddVideoCommentReply(cid int, author int, text string) {
	commentReply := VideoCommentReply{Cid: cid, Text: text, Author: author, Likes: 1}
	db.Insert(commentReply)
	return
}

func ClikckVideoEmoji(ctx *gin.Context) {
	session := sessions.Default(ctx)
	author := session.Get("userid")
	uauthor := int(author.(int64))
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
		DBaddVideoEmoji(vcid, uemoji, uauthor)
		return
	} else {
		var existedEmojiRecord VideoCommentEmojiRecord
		db.Where("author = ? and cid = ?", uauthor, cid).Get(&existedEmojiRecord)
		if existedEmojiRecord.Emoji == uemoji {
			DBdeleteVideoEmoji(vcid, uemoji, uauthor)
			return
		} else {
			DBchangeVideoEmoji(vcid, uemoji, existedEmojiRecord)
			return
		}
	}
}

func DBaddVideoEmoji(cid int, emoji int8, author int) {
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

func DBchangeVideoEmoji(cid int, emoji int8, emojiRecord VideoCommentEmojiRecord) {
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
	return
}

func DBdeleteVideoEmoji(cid int, emoji int8, author int) {
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
	return
}

func ClickVideoLike(ctx *gin.Context) {
	session := sessions.Default(ctx)
	author := session.Get("userid")
	uauthor := int(author.(int64))
	crid := ctx.PostForm("crid")
	vcrid, _ := strconv.Atoi(crid)
	count, _ := db.Where("author = ? and crid = ?", author, crid).Count(&VideoCommentReplyLikeRecord{})
	if count == 0 {
		DBaddVideoLike(uauthor, vcrid)
	} else {
		DBdeleteVideoLike(uauthor, vcrid)
	}
	return
}
func DBaddVideoLike(author int, crid int) {
	var commentReply VideoCommentReply
	db.ID(crid).Get(&commentReply)
	commentReply.Likes++
	commentReplyLikeRecord := VideoCommentReplyLikeRecord{Author: author, Crid: crid}
	db.Insert(&commentReplyLikeRecord)
	db.ID(crid).Update(&commentReply)
	return
}
func DBdeleteVideoLike(author int, crid int) {
	var commentReply VideoCommentReply
	db.ID(crid).Get(&commentReply)
	commentReply.Likes--
	db.Where("auhtor = ? and crid = ?", author, crid).Delete(&VideoCommentReplyLikeRecord{})
	db.ID(crid).Update(&commentReply)
	return
}

// 获取杂项数据
func GetVideoImg(c *gin.Context) {

	strid := c.Param("id")
	if strid == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	id, err := strconv.Atoi(strid)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //返回400
		return
	}
	var video Video
	_, err = db.ID(id).Get(&video)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error()) //500
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=cover.webp")
	c.Header("Vary", "Accept-Encoding")
	//c.Header("Content-Encoding", "br") //声明压缩格式,否则会被当作二进制文件下载
	c.File(video.CoverPath)
	return
}

func GetVideoTags(c *gin.Context) {
	var tagTexts []string
	var tagIds []int
	strid := c.Param("id")
	if strid == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	id, err := strconv.Atoi(strid)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //返回400 		return
	}
	var tags []Tag
	count, _ := db.Where("kind = ? AND pid = ?", 1, id).Count(&tags)
	if count == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	db.Where("kind = ? AND pid = ?", 1, id).Find(&tags)
	var tagModel TagModel
	var tid int
	for _, value := range tags {
		tid = int(value.Tid)
		db.ID(tid).Get(&tagModel)
		tagTexts = append(tagTexts, tagModel.Text)
		tagIds = append(tagIds, tid)
	}
	c.JSONP(http.StatusOK, gin.H{
		"tagtext": tagTexts,
		"tagid":   tagIds,
	})
	return
}

//两步上传

func UploadVideo(c *gin.Context) {
	formData, err := c.MultipartForm()
	files := formData.File["video"]
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	uuid := UUID.New()
	suuid := uuid.String()
	dst := "./temp/" + suuid
	var dstFile string
	var dstFiles []string
	os.MkdirAll(dst, os.ModePerm)
	for i, v := range files {
		dstFile = strings.Join([]string{dst, "/", strconv.Itoa(i), path.Ext(v.Filename)}, "")
		dstFiles = append(dstFiles, dstFile)
		if err1 := c.SaveUploadedFile(v, dstFile); err1 != nil {
			os.RemoveAll(dst)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	videoNeedtoCheck := VideoNeedToCheck{VideoPath: dstFiles}
	db.Insert(&videoNeedtoCheck)
	c.JSONP(http.StatusOK, gin.H{
		"id": videoNeedtoCheck.Id,
	})
	return
}

// TODO 先摸了

func SaveVideo(author int, src string, cscr string, title string, description string, uuid string) {
	var video Video
	video.Author = author
	video.Title = title
	video.Description = description
	video.Views = 0
	video.CoverPath = cscr
	//the error ffmpeg part
	err := ffmpeg.Input(src).Output(src+".mp4", ffmpeg.KwArgs{
		// "c:v": "libsvtav1",
	}).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		panic(err)
	}
	//
	// video.IpfsHash = Upload(src + ".mp4")
	db.Insert(&video)
	return
}

func DBaddVideoTag(vid int, tagid int) {
	tag := Tag{Tid: tagid, Pid: vid}
	db.Insert(tag)
	return
}
