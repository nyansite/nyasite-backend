package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"math"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/chai2010/webp"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	UUID "github.com/google/uuid"
)

func NewTag(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("is_login") != true {
		c.AbortWithStatus(http.StatusUnauthorized) //返回401
		return
	}
	tagname := c.PostForm("tagname")
	if has, _ := db.Exist(&TagModel{Text: tagname}); has == true {
		c.AbortWithStatus(StatusRepeatTag)
		return
	}
	db.Insert(&TagModel{Text: tagname})
	c.AbortWithStatus(http.StatusOK)
}

func AdminVideoPost(ctx *gin.Context) {
	vpg := ctx.Param("page")
	pg, err := strconv.Atoi(vpg)
	if err != nil || pg < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	var videos []Video
	var count int64 //总数,Count比rowsaffected更快(懒得用变量缓存了
	pg -= 1
	count, err = db.Count(&Video{})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) //500,正常情况下不会出现
		fmt.Println(err)
		return
	}
	db.Limit(20, pg*20).Find(&videos)
	ctx.JSON(http.StatusOK, gin.H{
		"Body":      videos,
		"PageCount": math.Ceil(float64(count) / 20), //总页数
	})
}

// 不审核直接上传,测试接口
func UploadVideo(c *gin.Context) {
	session := sessions.Default(c)
	title := c.PostForm("Title")
	description := c.PostForm("Description") //简介
	f, err1 := c.FormFile("file")
	cover, err2 := c.FormFile("cover")
	if err1 != nil || err2 != nil {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	author := session.Get("userid")
	vauthor := author.(int64)
	uauthor := uint(vauthor)
	uuid := UUID.New()
	sid := uuid.String()
	fpath := "./temporary/" + sid + path.Ext(f.Filename)
	cpath := "./temporary/" + sid + path.Ext(cover.Filename)
	c.SaveUploadedFile(f, fpath)
	c.SaveUploadedFile(cover, cpath)
	//transform cover into webp image
	fCover, _ := os.Open(cpath)
	image, _, _ := image.Decode(fCover)
	cpath = "./temporary/" + sid + ".webp"
	outfile, _ := os.Create(cpath)
	b := bufio.NewWriter(outfile)
	webp.Encode(b, image, &webp.Options{Lossless: false})
	//
	err1 = b.Flush()
	if err1 != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	go SaveVideo(uauthor, fpath, cpath, title, description, sid)
}

func GetSessionSecret() [][]byte {

	return nil
}
