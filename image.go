package main

import (
	"fmt"
	"image"
	"log"
	"net/http"

	_ "golang.org/x/image/webp"

	"github.com/gin-gonic/gin"
)

func PostImg(ctx *gin.Context) {
	file, err := ctx.FormFile("img")//这里已经读取文件了,文件过大会造成严重的性能问题
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	f, err := file.Open()
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	_, format, err := image.DecodeConfig(f)
	if err != nil {
		if err == image.ErrFormat { //只能识别加载Decode的类型
			fmt.Println(format)
		} else {
			panic(err)
		}
	}
	if format != "webp" { //不加载其他格式的话这个只能是webp
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)//415 不支持的媒体格式
		return
	}

	log.Println(format)
	// os.Create()
	//TODO saveimg
}
