package main

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type VerCode struct {
	Code    string
	GenTime int64
}

var VerCodeAllocMap map[string]VerCode

func AllocVerCodeResetPwd(c *gin.Context) {
	email := c.PostForm("email")
	IsUserExisted, _ := db.In("email", email).Exist(&User{})
	if !IsUserExisted {
		c.String(http.StatusBadRequest, "UserIsNotExisted")
	}
	verCode, isCodeExisted := VerCodeAllocMap[email]
	if isCodeExisted && (verCode.GenTime+60) >= time.Now().Unix() {
		timeInterval := verCode.GenTime + 60 - time.Now().Unix()
		c.String(http.StatusBadRequest, strconv.Itoa(int(timeInterval)))
		//lightsail server itself may send 429 ,so we use 400 instead
	} else {
		code := GenVerCode()
		VerCodeAllocMap[email] = VerCode{Code: code, GenTime: time.Now().Unix()}
		SendVerEmail(email, code, "重置密码")
		c.AbortWithStatus(http.StatusOK)
	}
}

func SendVerEmail(receiver string, verCode string, reason string) {

}

func VerifyVerCode(email string, code string) (bool, error) {
	if VerCodeAllocMap[email].Code == code {
		//if expired
		if (VerCodeAllocMap[email].GenTime + 300) >= time.Now().Unix() {
			delete(VerCodeAllocMap, email)
			return true, nil
		} else {
			expiredErr := errors.New("expired")
			delete(VerCodeAllocMap, email)
			return false, expiredErr
		}
	} else {
		incorrectErr := errors.New("incorrect")
		return false, incorrectErr
	}
}

func GenVerCode() string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, 6)
	for i := 0; i < 6; i++ {
		bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(bytes)
}
