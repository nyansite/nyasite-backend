package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"crypto/rand"
	"crypto/sha512"
	_ "encoding/json"
	_ "fmt"
	"regexp"
	"strconv"
)

func GetSelfUserData(c *gin.Context) {

	session := sessions.Default(c)
	var user User

	if session.Get("is_login") == true {
		userid := session.Get("userid")
		level := session.Get("level")
		db.First(&user, userid)
		mail := user.Email
		c.JSON(http.StatusOK, gin.H{
			"userid": userid,
			"mail":   mail,
			"level":  level,
		})
	} else {
		c.AbortWithStatus(http.StatusUnauthorized) //返回401
		return
	}
}

func GetUserData(c *gin.Context) {
	id := c.Param("id")
	nid, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized) //返回401
		return
	}
	var user User
	db.First(&user, nid)
	c.JSON(http.StatusOK, gin.H{
		"name":  user.Name,
		"level": user.Level,
	})
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("is_login") == true {
		c.AbortWithStatus(StatusAlreadyLogin)
		return
	}
	username, passwd := c.PostForm("username"), c.PostForm("passwd") //传入的用户名也有可能是邮箱
	if username == "" || passwd == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	var user User
	if db.First(&user, "Name = ? OR Email = ?", username, username).RowsAffected == 0 { //用户不存在
		c.AbortWithStatus(StatusUserNameNotExist)
		return
	}
	if !check_passwd(user.Passwd, []byte(passwd)) {
		c.AbortWithStatus(StatusPasswordError)
		return
	}

	session.Set("userid", user.ID)
	session.Set("is_login", true)
	session.Set("level", user.Level)
	session.Save()
	c.AbortWithStatus(StatusLoginOK)
}

func Register(c *gin.Context) {
	username, passwd, mail := c.PostForm("username"), c.PostForm("passwd"), c.PostForm("email")

	if username == "" || passwd == "" || mail == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}

	if reg := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`); !reg.MatchString(mail) { //检测前端也要做一遍
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//上面判断输入是否合法,下面判断用户是否已经存在

	if db.First(&User{}, "Name = ?", username).RowsAffected != 0 {
		c.AbortWithStatus(StatusRepeatUserName)
		return
	}
	if db.First(&User{}, "Email = ?", mail).RowsAffected != 0 {
		c.AbortWithStatus(StatusRepeatEmail)
		return
	}

	user := User{Name: username, Passwd: encrypt_passwd([]byte(passwd)), Email: mail}
	db.Create(&user)
	c.AbortWithStatus(StatusUserCreatedOK)
}

func encrypt_passwd(passwd []byte) []byte { //加密密码,带盐
	salte, _ := rand.Prime(rand.Reader, 64) //普普通通的64位盐,8字节
	salt := salte.Bytes()

	passwd_sha := sha512.Sum512(passwd)          //密码的sha
	saltpasswd := append(passwd_sha[:], salt...) //加盐
	safe_passwd := sha512.Sum512_256(saltpasswd) //这一步才算加密,512/256指生成512之后截断成256,安全性一样

	return append(safe_passwd[:], salt...) //保存盐
}

func check_passwd(passwd []byte, passwd2 []byte) bool {
	//获取盐
	salt := passwd[32:]
	passwd = passwd[:32]

	passwd2_sha := sha512.Sum512(passwd2)
	saltpasswd2 := append(passwd2_sha[:], salt...)
	safe_passwd := sha512.Sum512_256(saltpasswd2)

	ret := true
	for i, v := range passwd {
		if v != safe_passwd[i] {
			ret = false
			//不要break防止时间攻击
		}
	}
	return ret
}
