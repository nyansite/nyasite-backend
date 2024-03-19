package main_test

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func LimitRequestBody(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		err := c.Request.ParseMultipartForm(maxSize)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return
		}
		c.Next()
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	// 定义中间件函数，用于设置不同路由的请求体大小限制
	r.POST("/upload", LimitRequestBody(1<<20),func(ctx *gin.Context) {
		f, err := ctx.FormFile("file")
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
        log.Println(f.Size)

	})
	return r
}

func TestFileSizeLimit(t *testing.T) {
	// 设置路由
	router := setupRouter()

	// 创建一个超出大小限制的文件
	data := make([]byte, 1<<20)
	file := bytes.NewReader(data)

	// 创建一个带有 multipart 数据的请求
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	part, _ := writer.CreateFormFile("file", "example.txt")
	io.Copy(part, file)
	writer.Close()
    log.Println("start")
	// 创建一个 POST 请求到 "/upload" 路由，携带 multipart 数据
	req := httptest.NewRequest("POST", "/upload", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// 使用路由的 ServeHTTP 方法来执行请求
	router.ServeHTTP(w, req)


	println(w.Code)
	println(w.Body.String())
}
