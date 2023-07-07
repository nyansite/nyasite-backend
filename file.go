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

func Upload(path string) string {
	sh := shell.NewShell("localhost:5001") //根据服务器ipfs客户端的配置改
	uploadfile, _ := os.OpenFile(path, os.O_RDONLY, 777)
	cid, err := sh.Add(uploadfile)
	if err != nil {
		println(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	return cid
}
