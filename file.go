package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/andybalholm/brotli"
	"io"
	"net/http"
	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	// "gorm.io/gorm"
	"errors"
	"strconv"
)

var (
	NotFound     = errors.New("梦里啥都有,这里什么都没有")
	NoIpfsDaemon = errors.New("ipfs daemon被你吃了?")
)

func BrowseVideo(ctx *gin.Context) {

	vid := ctx.Param("vid")
	if vid == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, "/1")
		return
	}

	id, err := strconv.Atoi(vid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	if id < 1 {
		ctx.String(http.StatusBadRequest, "你搁这翻空气呢?")
		return
	}
	var videos []Video
	var count int64 //总数,Count比rowsaffected更快(懒得用变量缓存了
	id -= 1
	db.Model(&Video{}).Count(&count)
	db.Limit(20).Offset(id * 20).Find(&videos)
	// for i, v:= range videos{

	// }
	ctx.HTML(http.StatusOK, "browsevideo.html", gin.H{
		"list": videos,
	})
}

// TODO
// 从ipfs中查看文件列表/文件名
// func BrowseFiles(ctx *gin.Context) {
// 	sh := shell.NewLocalShell() //需要挂着ipfs daemon
// 	if sh == nil {
// 		ctx.AbortWithStatus(http.StatusInternalServerError) //500
// 		return
// 	}
// 	sh.SetTimeout(1145140000) //为啥单位是纳秒???

// 	ctx2 := context.Background() //直接传给ipfs daemon,密码之类的,(很没用
// 	// sh.FilesMkdir(ctx2, "/img")		//防小天才

// 	fls, err := sh.FilesLs(ctx2, "")
// 	if err != nil {
// 		ctx.AbortWithStatus(http.StatusInternalServerError) //500
// 		fmt.Println(err)
// 		return
// 	}
// 	for _, v := range fls{
// 		_ = v
// 	}
// }

// 从ipfs获取文件,测试用
// 只有用AddFile上传的文件才能用,因为存储的是压缩数据
func GetFile(path string) ([]byte, error) {
	sh := shell.NewLocalShell() //需要挂着ipfs daemon
	if sh == nil {
		return nil, NoIpfsDaemon
	}
	sh.SetTimeout(1145140000) //为啥单位是纳秒???
	brfr, _ := sh.Cat(path)
	if brfr == nil {

		return nil, NotFound
	}
	// cl := brotli.NewReader(brfr)
	// buf, _ := io.ReadAll(cl)
	buf2, _ := io.ReadAll(brfr)

	// ctx.String(http.StatusOK, string(buf))
	return buf2, nil

}

func AddFile(r io.Reader, path string) error {
	buf := bytes.Buffer{} //输出的缓冲区
	buf2 := bytes.Buffer{}
	buf2.ReadFrom(r)

	cl := brotli.NewWriter(&buf)
	cl.Write(buf2.Bytes())
	cl.Close()

	sh := shell.NewLocalShell() //需要挂着ipfs daemon
	if sh == nil {
		return NoIpfsDaemon
	}

	ctx := context.Background()
	err := sh.FilesWrite(ctx, path, &buf)
	return err	//正常情况下应该是nil(大概)
}

func AddFileT(c *gin.Context) {
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
