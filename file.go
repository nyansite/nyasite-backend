package main

import (
	// "fmt"
	"net/http"
	"io"
	"github.com/gin-gonic/gin"
	"bytes"
	// "github.com/ipfs/boxo/files"
)

// 从ipfs获取文件
func GetFile(c *gin.Context) {

}

func AddFile(c *gin.Context) {
	//TODO 权限

	// files.NewBytesFile()
	
	f, _, err := c.Request.FormFile("file")    
	defer f.Close()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	// files.NewBytesFile(buf.Bytes())
}
