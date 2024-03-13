package main

import (
	"bufio"
	"fmt"
	_ "golang.org/x/image/webp"
	"image"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func PostImg(ctx *gin.Context) {
	img, err := ctx.FormFile("img")
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	f, err := img.Open()
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	_, format, err := image.DecodeConfig(bufio.NewReader(f))
	if err != nil {
		if err == image.ErrFormat {//只能识别加载Decode的类型
			fmt.Println(format)
		} else {
			panic(err)
		}
	}
	if format == "webp"{//不加载其他格式的话这个只能是webp
		//TODO saveimg
	}
}
func test() {
	file, _ := os.Open("/home/tksmly/Downloads/116810783_p0.webp")
	_, format, err := image.DecodeConfig(bufio.NewReader(file))
	if err != nil {
		if err == image.ErrFormat {
			fmt.Println(format)
		} else {
			panic(err)
		}
	}
	fmt.Println(format)
}
