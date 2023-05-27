package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
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

	sh := shell.NewLocalShell()	//需要挂着ipfs daemon
	sh.Add(f)
	fmt.Println()
}
