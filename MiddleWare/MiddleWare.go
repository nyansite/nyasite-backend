//修改自 github.com/anargu/gin-brotli (MIT)
//并不需要压缩,数据库存储已压缩文件
package middleware

import (
	"net/http"
	"strings"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/andybalholm/brotli"

)
//Write需要的接口↓
func (br *brotliWriter) WriteString(s string) (int, error) {
	return br.writer.Write([]byte(s))
}

func (br *brotliWriter) Write(data []byte) (int, error) {
	return br.writer.Write(data)
}

func (br *brotliWriter) WriteHeader(code int) {
	br.Header().Del("Content-Length")
	br.ResponseWriter.WriteHeader(code)
}
//Write需要的接口↑

type brotliWriter struct {
	gin.ResponseWriter
	writer *brotli.Writer
}
func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "br") ||
		strings.Contains(req.Header.Get("Connection"), "Upgrade") ||
		strings.Contains(req.Header.Get("Content-Type"), "text/event-stream") {

		return false
	}

	extension := filepath.Ext(req.URL.Path)
	if len(extension) < 4 { // fast path
		return true
	}

	// if skip := containsString([]string{".png", ".gif", ".jpeg", ".jpg", ".mp3", ".mp4"}, extension); skip {
	// 	return false
	// } else {
	// 	return true
	// }
	return true
}

func Br() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !shouldCompress(ctx.Request){
			ctx.String(http.StatusUpgradeRequired, "请使用支持brotli的浏览器(2016年之后)")
			ctx.Abort()
			return
		}
		ctx.Header("Content-Encoding", "br")
		ctx.Header("Vary", "Accept-Encoding")
		brWriter := brotli.NewWriterOptions(ctx.Writer, brotli.WriterOptions{
			Quality: 4,
			LGWin:   11,
		})
		
		ctx.Writer = &brotliWriter{ctx.Writer, brWriter}

		defer func() {
			brWriter.Close()
			ctx.Header("Content-Length", fmt.Sprint(ctx.Writer.Size()))
		}()
		ctx.Next()
	}
}
