package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPICUItoken(c *gin.Context) {
	const url = "https://picui.cn/api/v1/images/tokens"
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
	//PICUItoken from token.go
	req.Header.Set("Authorization", "Bearer "+"PICUItoken")
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
	b, err := ioutil.ReadAll(resp.Body)
	println(string(b))
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
