package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"net/http"
)

// 从ipfs获取文件请使用ipfs网关

func AddFile(c *gin.Context) {
	//TODO 权限

	session := sessions.Default(c)
	if session.Get("is_login") != true {
		c.AbortWithStatus(http.StatusUnauthorized) //返回401
		return
	}

	f, _, err := c.Request.FormFile("file")
	defer f.Close()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}

	sh := shell.NewLocalShell() //需要挂着ipfs daemon
	sh.Add(f)
	fmt.Println()
}
