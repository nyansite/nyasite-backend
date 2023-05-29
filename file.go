package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	// "io"
	// "net/http"

	"github.com/andybalholm/brotli"
	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
)

// 从ipfs获取文件,测试用
// 只有用AddFile上传的文件才能用,因为存储的是压缩数据
func GetFile(ctx *gin.Context) {
	ctx.Header("Content-Encoding", "br")  //声明压缩格式,否则会被当作二进制文件下载
	ctx.Header("Vary", "Accept-Encoding") //客户端使用缓存

	head := ctx.Query("file")
	if head == "" {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	sh := shell.NewLocalShell() //需要挂着ipfs daemon
	if sh == nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) //500
		return
	}
	sh.SetTimeout(300000000) //为啥单位是纳秒???
	brfr, err := sh.Cat(head)
	if brfr == nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusNotFound) //404
		return
	}
	// cl := brotli.NewReader(brfr)
	// buf, _ := io.ReadAll(cl)
	buf2, _ := io.ReadAll(brfr)

	// ctx.String(http.StatusOK, string(buf))
	ctx.Data(http.StatusOK, "text/plain", buf2)
}

func AddFile(c *gin.Context) {
	//TODO 权限
	//先注释掉因为要测试
	// session := sessions.Default(c)
	// if session.Get("is_login") != true {
	// 	c.AbortWithStatus(http.StatusUnauthorized) //返回401
	// 	return
	// }

	f, err := c.FormFile("file")

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	ff, _ := f.Open()
	defer ff.Close()

	buf := bytes.Buffer{} //输出的缓冲区
	b2 := bytes.Buffer{}  //输入的缓冲区,因为multipart.File未实现io.Writer
	b2.ReadFrom(ff)

	cl := brotli.NewWriter(&buf)
	cl.Write(b2.Bytes())
	cl.Close()

	sh := shell.NewLocalShell() //需要挂着ipfs daemon
	if sh == nil {
		c.AbortWithStatus(http.StatusInternalServerError) //500
		return
	}
	path, _ := sh.Add(&buf)
	fmt.Println(path)
	c.String(http.StatusOK, path)
}
