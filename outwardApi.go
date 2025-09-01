package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 申请图床的临时token
func GetWMINGtoken(c *gin.Context) {
	const url = "https://wmimg.com/api/v1/images/tokens"
	//申请一个token,有效时间10个小时
	formDataMap := map[string]int{
		"num":     1,
		"seconds": 36000,
	}
	formDataBytes, _ := json.Marshal(formDataMap)
	formDataBytesReader := bytes.NewReader(formDataBytes)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, formDataBytesReader)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//WMINGtoken from token.go
	req.Header.Set("Authorization", "Bearer "+WMINGtoken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Length", "42")
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		c.AbortWithStatus(http.StatusBadGateway)
	}
	var tokenJson map[string]interface{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err1 := json.Unmarshal(b, &tokenJson)
	if err1 != nil {
		c.AbortWithError(http.StatusInternalServerError, err1)
		return
	}
	token := tokenJson["data"].(map[string]interface{})["tokens"].([]interface{})[0].(map[string]interface{})["token"].(string)
	c.String(http.StatusOK, "%v", token)
	return
}

// 作为tus上传服务的中介
// https://developers.cloudflare.com/stream/uploading-videos/direct-creator-uploads/#advanced-upload-flow-using-tus-for-large-videos
func UploadVideoUrl(c *gin.Context) {
	length := c.Request.Header.Get("Upload-Length")
	metadata := c.Request.Header.Get("Upload-Metadata")
	const endpoint = "https://api.cloudflare.com/client/v4/accounts/" + CloudflareAccount + "/stream?direct_user=true"
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	req.Header.Set("Authorization", "bearer "+CloudflareStreamToken)
	req.Header.Set("Tus-Resumable", "1.0.0")
	req.Header.Set("Upload-Length", length)
	req.Header.Set("Upload-Metadata", metadata)
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		c.AbortWithStatus(http.StatusBadGateway)
	}
	c.Header("Location", resp.Header.Get("location"))
	c.Header("Access-Control-Expose-Headers", "Location")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Origin", "*")
	c.String(http.StatusOK, "%v", resp.Header.Get("stream-media-id"))
}

// 获取视频hls播放链接
func GetLinkPlaybackHls(c *gin.Context) {
	//当请求错误时会有panic错误导致整个程序停止，即使cloudflare也不能保证稳定访问（悲）
	defer func() {
		if err := recover(); err != nil {
			println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	uid := c.Param("uid")
	endpoint := "https://api.cloudflare.com/client/v4/accounts/" + CloudflareAccount + "/stream/" + uid
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	req.Header.Set("Authorization", "bearer "+CloudflareStreamToken)
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		c.AbortWithStatus(http.StatusBadGateway)
	}
	var linkJson map[string]interface{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	err1 := json.Unmarshal(b, &linkJson)
	if err1 != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	link := linkJson["result"].(map[string]interface{})["playback"].(map[string]interface{})["hls"].(string)
	c.String(http.StatusOK, "%v", link)

}
