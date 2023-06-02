package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"github.com/andybalholm/brotli"
	shell "github.com/ipfs/go-ipfs-api"
)

var (
	NotFound     = errors.New("梦里啥都有,这里什么都没有")
	NoIpfsDaemon = errors.New("ipfs daemon被你吃了?")
)

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
	sh.FilesFlush(ctx, "")
	return err //正常情况下应该是nil(大概)

}

func Addpath(src string, dst string) error {
	sh := shell.NewLocalShell() //需要挂着ipfs daemon
	if sh == nil {
		return NoIpfsDaemon
	}
	ctx := context.Background()

	path, err := sh.AddDir(src, shell.Pin(false))
	if err != nil {
		return err
	}
	sh.FilesMkdir(ctx, "/video")	//没啥用,但是防小天才
	err = sh.FilesCp(ctx, "/ipfs/"+path, dst)
	if err != nil {
		return err
	}
	
	sh.FilesFlush(ctx, "")
	return err
}
