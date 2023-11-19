package main

import (
	"fmt"
	"net/http"

	"crypto/rand"
	"crypto/sha512"
	"regexp"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSelfUserData(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("is_login") == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var user User
	userid := session.Get("userid")
	level := session.Get("level")
	vuserid, _ := userid.(int64)
	db.ID(vuserid).Get(&user)
	mail := user.Email

	session.Flashes() //重新set cookie,使得cookie生命周期重置,但是值不会重置
	session.Save()
	println(vuserid)
	c.JSON(http.StatusOK, gin.H{
		"name":   user.Name,
		"userid": userid,
		"mail":   mail,
		"level":  level,
		"avatar": DBGetUserDataShow(int(vuserid)).Avatar,
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
		"name":  user.Name,
		"level": user.Level,
	})
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("is_login") != nil {
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
	if !check_passwd(user.Passwd, []byte(passwd)) {
		c.AbortWithStatus(StatusPasswordError)
		return
	}
	session.Set("userid", user.Id)
	session.Set("is_login", true)
	session.Set("level", user.Level)
	session.Set("pwd-8", user.Passwd[:8]) //更改密码后其他已登录设备会退出
	session.Save()

	c.AbortWithStatus(http.StatusOK)
}

func Register(c *gin.Context) {
	username, passwd, mail, avatar := c.PostForm("username"), c.PostForm("passwd"), c.PostForm("email"), c.PostForm("avatar")

	if username == "" || passwd == "" || mail == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}

	if reg := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`); !reg.MatchString(mail) { //检测前端也要做一遍
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

	user := User{Name: username, Passwd: encrypt_passwd([]byte(passwd)), Email: mail, Avatar: avatar}
	_, err := db.Insert(&user)
	if err != nil {
		fmt.Println(err)
	}
	c.AbortWithStatus(http.StatusOK)
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
			//不要break防止时间攻击(也许不需要)
		}
	}
	return ret
}

func DBGetUserDataShow(userid int) UserDataShow {
	var userDataShow UserDataShow
	var user User
	db.ID(userid).Get(&user)
	userDataShow.Name = user.Name
	userDataShow.Avatar = user.Avatar
	userDataShow.Id = user.Id
	return userDataShow
}
