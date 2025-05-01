package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/gmail/v1"
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

func AllocVerCodeRegister(c *gin.Context) {
	email := c.PostForm("email")
	verCode, isCodeExisted := VerCodeAllocMap[email]
	if isCodeExisted && (verCode.GenTime+60) >= time.Now().Unix() {
		timeInterval := verCode.GenTime + 60 - time.Now().Unix()
		c.String(http.StatusBadRequest, strconv.Itoa(int(timeInterval)))
		//lightsail server itself may send 429 ,so we use 400 instead
	} else {
		code := GenVerCode()
		VerCodeAllocMap[email] = VerCode{Code: code, GenTime: time.Now().Unix()}
		SendVerEmail(email, code, "注册账户")
		c.AbortWithStatus(http.StatusOK)
	}
}

func GenMail(requestType string, OTP string) string {
	var builder strings.Builder
	builder.WriteString("<html><head><meta charset='UTF-8' /><style>body {margin: 0;padding: 0;color: #333;background-color: #fff;}.container {margin: 0 auto;width: 100%;max-width: 600px;padding: 0 0px;padding-bottom: 10px;border-radius: 5px;line-height: 1.8;}.header {border-bottom: 1px solid #eee;}.header a {font-size: 1.4em;color: #000;text-decoration: none;font-weight: 600;}.content {min-width: 700px;overflow: auto;line-height: 2;}.otp {background: linear-gradient(to right, #00bc69 0, #00bc88 50%, #00bca8 100%);margin: 0 auto;width: max-content;padding: 0 10px;color: #fff;border-radius: 4px;}.footer {color: #aaa;font-size: 0.8em;line-height: 1;font-weight: 300;}.email-info {color: #666666;font-weight: 400;font-size: 13px;line-height: 18px;padding-bottom: 6px;}.email-info a {text-decoration: none;color: #00bc69;}</style></head><body><div class='container'><div class='header'><a>验证您的喵站帐号</a></div><br /><p>我们从您的喵站帐号接收到了一个")
	builder.WriteString(requestType)
	builder.WriteString("请求。为了安全, 请输入一下一次性验证码来验证您的身份<br /><b>您的一次性验证码是:</b></p><h2 class='otp'>")
	builder.WriteString(OTP)
	builder.WriteString("</h2><p style='font-size: 0.9em'><strong>一次性验证码的有效期为5分钟</strong><br /><br />如果您没有")
	builder.WriteString(requestType)
	builder.WriteString("，请忽略这封电子邮件。请确保你的一次性验证码不会泄漏给任何人<br /><strong>不要把这封电子邮件转发给任何人</strong><br /><br /><strong>感谢您选择喵站</strong><br /></p><hr style='border: none; border-top: 0.5px solid #131111' /><div class='footer'><p>本邮件为自动发送，请不要尝试回复</p><!-- <p>如果想知道更多喵站和您的帐号的信息, visit<strong>[Name]</strong></p> --></div></div></body></html>")
	return builder.String()
}

func SendVerEmail(receiver string, verCode string, reason string) {
	println(verCode)
	var message gmail.Message
	messageStr := []byte(
		"From: nyasite@gmail.com\r\n" +
			"To: " + receiver + "\r\n" +
			"Content-Type: text/html;\r\n" +
			"Subject: Verified Code\r\n\r\n" + GenMail(reason, verCode))
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)
	_, err := emailSrv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Message sent!")
	}
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
