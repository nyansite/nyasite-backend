package main

import (
	"fmt"
	"net/http"
	"time"

	"crypto/rand"
	"crypto/sha512"
	"regexp"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func reloadJWT(user User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"email":    user.Email,
		"username": user.Name,
		"picture":  user.Avatar,
		"exp":      time.Now().Add(30 * 24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("nyasite"))
	return tokenString
}

func GetSelfUserData(c *gin.Context) {
	session := sessions.Default(c)
	is_login, _ := c.Cookie("is_login")
	if is_login != "true" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var user User
	userid := session.Get("userid")

	vuserid, _ := userid.(int64)
	db.ID(vuserid).Get(&user)
	mail := user.Email
	if (int((int(time.Now().Unix())+user.Timezone)/86400) - int((user.LTC+user.Timezone)/86400)) >= 1 {
		user.Level = user.Level + 1
		user.LTC = int(time.Now().Unix())
		db.ID(vuserid).Update(&user)
	}
	session.Flashes() //重新set cookie,使得cookie生命周期重置,但是值不会重置
	session.Save()
	//刷新jwt
	tokenString := reloadJWT(user)
	c.SetCookie("token", tokenString, 1200000, "/", "", true, true)
	c.SetCookie("is_login", "true", 1200000, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"name":     user.Name,
		"userid":   userid,
		"mail":     mail,
		"level":    user.Level,
		"avatar":   user.Avatar,
		"timezone": user.Timezone,
	})
}

func GetUserData(c *gin.Context) {
	id := c.Param("id")
	nid, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest) //返回400
		return
	}
	var user User
	db.ID(nid).Get(&user)

	c.JSON(http.StatusOK, gin.H{
		"name":   user.Name,
		"level":  user.Level,
		"avatar": user.Avatar,
	})
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	is_login, _ := c.Cookie("is_login")
	if is_login == "true" {
		c.AbortWithStatus(StatusAlreadyLogin)
		return
	}
	username, passwd := c.PostForm("username"), c.PostForm("passwd") //传入的用户名也有可能是邮箱
	if username == "" || passwd == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	var user User
	if has, _ := db.Where("Name = ? OR Email = ?", username, username).Get(&user); has == false { //用户不存在
		c.Status(StatusUserNameNotExist)
		return
	}
	if !check_passwd(user.Passwd, passwd) {
		c.AbortWithStatus(StatusPasswordError)
		return
	}
	session.Set("userid", user.Id)
	session.Set("pwd-8", user.Passwd[:8]) //更改密码后其他已登录设备会退出
	session.Save()
	//刷新jwt
	tokenString := reloadJWT(user)
	c.SetCookie("token", tokenString, 1200000, "/", "", true, true)
	c.SetCookie("is_login", "true", 1200000, "/", "", true, true)
	c.AbortWithStatus(http.StatusOK)
}

var regCompile = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)

func Register(c *gin.Context) {
	username, passwd, mail := c.PostForm("username"), c.PostForm("passwd"), c.PostForm("email")

	if username == "" || passwd == "" || mail == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	if !regCompile.MatchString(mail) { //检测前端也要做一遍
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//上面判断输入是否合法,下面判断用户是否已经存在

	if has, _ := db.Exist(&User{Name: username}); has == true {
		c.AbortWithStatus(StatusRepeatUserName)
		return
	}
	if has, _ := db.Exist(&User{Email: mail}); has == true {
		c.AbortWithStatus(StatusRepeatEmail)
		return
	}
	timezone := c.PostForm("timezone")
	timezoneOffSet, _ := strconv.Atoi(timezone)
	user := User{Name: username, Passwd: encrypt_passwd(passwd), Email: mail,
		Timezone: timezoneOffSet, Avatar: ("https://ui-avatars.com/api/?background=b3c6d7&name=" + username)}
	_, err := db.Insert(&user)
	if err != nil {
		fmt.Println(err)
	}
	c.AbortWithStatus(http.StatusOK)
}

func encrypt_passwd(passwds string) []byte { //加密密码,带盐
	salte, _ := rand.Prime(rand.Reader, 64) //普普通通的64位盐,8字节
	salt := salte.Bytes()
	passwd := []byte(passwds)
	passwd_sha := sha512.Sum512(passwd)          //密码的sha
	saltpasswd := append(passwd_sha[:], salt...) //加盐
	safe_passwd := sha512.Sum512_256(saltpasswd) //这一步才算加密,512/256指生成512之后截断成256,安全性一样

	return append(safe_passwd[:], salt...) //保存盐
}

func check_passwd(passwd []byte, passwd2s string) bool {
	//获取盐
	salt := passwd[32:]
	passwd = passwd[:32]
	passwd2 := []byte(passwd2s)

	passwd2_sha := sha512.Sum512(passwd2)
	saltpasswd2 := append(passwd2_sha[:], salt...)
	safe_passwd := sha512.Sum512_256(saltpasswd2)

	ret := true
	for i, v := range passwd {
		if v != safe_passwd[i] {
			ret = false
			//不要break防止时间攻击(也许不需要)
		}
	}
	return ret
}

//获取用户id

func GetUserIdWithCheck(c *gin.Context) int {
	is_login, _ := c.Cookie("is_login")

	if is_login == "true" {
		session := sessions.Default(c)
		author := session.Get("userid")
		return int(author.(int64))
	} else {
		return -1
	}
}

func GetUserIdWithoutCheck(c *gin.Context) int {
	session := sessions.Default(c)
	author := session.Get("userid")
	return int(author.(int64))
}

//获取用户信息(程序内)

func DBGetUserDataShow(userid int) UserDataShow {
	var user User
	db.ID(userid).Get(&user)
	userDataShow := UserDataShow{
		Id:     user.Id,
		Name:   user.Name,
		Avatar: user.Avatar,
	}
	return userDataShow
}

//更改用户信息

func ChangeAvatar(c *gin.Context) {
	uauthor := GetUserIdWithoutCheck(c)
	avatar := c.PostForm("avatar")
	var user User
	db.ID(uauthor).Get(&user)
	user.Avatar = avatar
	db.ID(uauthor).Update(&user)
	return
}

func ChangeName(c *gin.Context) {
	uauthor := GetUserIdWithoutCheck(c)
	name := c.PostForm("name")
	var user User
	db.ID(uauthor).Get(&user)
	user.Name = name
	user.Level = user.Level - 8
	db.ID(uauthor).Update(&user)
	return
}

func ChangeTimeZone(c *gin.Context) {
	uauthor := GetUserIdWithoutCheck(c)
	timezone := c.PostForm("timezone")
	timezoneOffSet, _ := strconv.Atoi(timezone)
	var user User
	db.ID(uauthor).Get(&user)
	user.Timezone = timezoneOffSet
	db.ID(uauthor).Update(&user)
	return
}
