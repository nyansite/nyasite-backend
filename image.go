package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	_ "golang.org/x/image/webp"
	"image"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func PostImg(ctx *gin.Context) {
	file, err := ctx.FormFile("img") //这里已经读取文件了,文件过大会造成严重的性能问题
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	f, err := file.Open()
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	_, format, err := image.DecodeConfig(f)
	if err != nil { //只能识别加载Decode的类型
		log.Println(err)
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType) //415 不支持的媒体格式
		return
	}
	if format != "webp" { //不加载其他格式的话这个只能是webp
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType) //415 不支持的媒体格式
		return
	}
	SaveImg(f)
	//TODO saveimg
}

func SaveImg(file multipart.File) error {
	var ff []byte
	file.Read(ff)
	hash := sha256.Sum256(ff)
	_, err := os.OpenFile("img/"+hex.EncodeToString(hash[:])+".webp", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err){
			log.Println(err)
			return err
		}
		log.Panic(err)
	}
	return nil
}
