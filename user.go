package main

import (
	"fmt"
	"net/http"
	"time"

	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func reloadJWT(user User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"email":    user.Email,
		"username": user.OriginName,
		"fullname": user.Name,
		"picture":  user.Avatar,
		"exp":      time.Now().Add(30 * 24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(TokenJWTSecret))
	return tokenString
}

func QuitLogin(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", true, true)
	c.SetCookie("is_login", "false", -1, "/", "", true, true)
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
	db.ID(int(userid.(int64))).Get(&user)
	mail := user.Email
	c.JSON(http.StatusOK, gin.H{
		"name":   user.Name,
		"userid": userid,
		"mail":   mail,
		"level":  user.Level,
		"avatar": user.Avatar,
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
		c.AbortWithStatus(http.StatusForbidden) //前端应该跳转,而不是重复请求登入
		return
	}
	username, passwd := c.PostForm("username"), c.PostForm("passwd") //传入的用户名也有可能是邮箱
	if username == "" || passwd == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	var user User
	if has, _ := db.Where("Name = ? OR Email = ?", username, username).Get(&user); !has { //用户不存在
		c.AbortWithStatus(http.StatusUnauthorized) //401 不存在或错误的身份验证
		return
	}
	if !check_passwd(user.Passwd, passwd) {
		c.AbortWithStatus(http.StatusUnauthorized) //401
		return
	}
	session.Set("userid", user.Id)
	session.Set("pwd-8", user.Passwd[:8]) //高8位,更改密码后其他已登录设备会退出
	session.Save()
	//刷新jwt
	tokenString := reloadJWT(user)
	c.SetCookie("token", tokenString, 1200000, "/", "", true, true)
	c.SetCookie("is_login", "true", 1200000, "/", "", true, true)
	c.AbortWithStatus(http.StatusOK)
}

var regCompile = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)

func Register(c *gin.Context) {
	username, passwd, email, code := c.PostForm("username"), c.PostForm("passwd"), c.PostForm("email"), c.PostForm("verCode")

	if username == "" || passwd == "" || email == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	if !regCompile.MatchString(email) { //检测前端也要做一遍
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//上面判断输入是否合法,下面判断用户是否已经存在

	if has, _ := db.Exist(&User{Name: username}); has {
		c.String(http.StatusConflict, "NameUsed")
		c.Abort()
		return
	}
	if has, _ := db.Exist(&User{Email: email}); has {
		c.String(http.StatusConflict, "EmailAddressUsed")
		c.Abort()
		return
	}
	canReset, err := VerifyVerCode(email, code)
	if canReset {
		user := User{Name: username, Passwd: encrypt_passwd(passwd), Email: email, OriginName: username}
		_, err := db.Insert(&user)
		if err != nil {
			fmt.Println(err)
		}
		c.AbortWithStatus(http.StatusOK)
	} else if err.Error() == "expired" {
		c.String(http.StatusBadRequest, "expired")
	} else if err.Error() == "incorrect" {
		c.String(http.StatusBadRequest, "incorrectVerCode")
	}

}

func ResetPwd(c *gin.Context) {
	var user User
	newPwd := c.PostForm("pwd")
	email := c.PostForm("email")
	code := c.PostForm("verCode")
	db.In("email", email).Get(&user)
	canReset, err := VerifyVerCode(email, code)
	if canReset {
		user.Passwd = encrypt_passwd(newPwd)
		db.In("email", email).Update(&user)
	} else if err.Error() == "expired" {
		c.String(http.StatusBadRequest, "expired")
	} else if err.Error() == "incorrect" {
		c.String(http.StatusBadRequest, "incorrectVerCode")
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

// 这玩意应该放在别的地方
func Refresh(c *gin.Context) {
	session := sessions.Default(c)
	userid := session.Get("userid")
	is_login, _ := c.Cookie("is_login")
	if is_login != "true" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var user User
	if has, _ := db.ID(int(userid.(int64))).Get(&user); !has { //用户不存在
		c.SetCookie("token", "", -1, "/", "", true, true)
		c.SetCookie("is_login", "false", -1, "/", "", true, true)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if string(session.Get("pwd-8").([]byte)) != string(user.Passwd[:8]) {
		c.SetCookie("token", "", -1, "/", "", true, true)
		c.SetCookie("is_login", "false", -1, "/", "", true, true)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	session.Flashes() //重新set cookie,使得cookie生命周期重置,但是值不会重置
	session.Save()
	//刷新jwt
	tokenString := reloadJWT(user)
	c.SetCookie("token", tokenString, 1200000, "/", "", true, true)
	c.SetCookie("is_login", "true", 1200000, "/", "", true, true)
}

func ClockIn(c *gin.Context) {
	userid := GetUserIdWithoutCheck(c)
	timezoneStr := c.PostForm("timezone")
	timezone, _ := strconv.Atoi(timezoneStr)
	var user User
	db.ID(userid).Get(&user)
	//不在同一天不是一个整数
	if (int((int(time.Now().Unix())+timezone)/86400) - int((user.LTC+timezone)/86400)) >= 1 {
		user.Level = user.Level + 1
		user.LTC = int(time.Now().Unix())
		db.ID(userid).Update(&user)
	}
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
	salt := passwd[32:]  //32位开始到结束
	passwd = passwd[:32] //开始到32位结束
	passwd2 := []byte(passwd2s)

	passwd2_sha := sha512.Sum512(passwd2)
	saltpasswd2 := append(passwd2_sha[:], salt...)
	safe_passwd := sha512.Sum512_256(saltpasswd2)

	//防止时间差攻击
	return subtle.ConstantTimeCompare(passwd, safe_passwd[:]) == 1 //甜蜜的go没法直接int>bool
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
	userDisplay := UserDataShow{
		Id:     user.Id,
		Name:   user.Name,
		Avatar: user.Avatar,
	}
	return userDisplay
}

//更改用户信息

func ChangeAvatar(c *gin.Context) {
	uauthor := GetUserIdWithoutCheck(c)
	avatar := c.PostForm("avatar")
	var user User
	db.ID(uauthor).Get(&user)
	user.Avatar = avatar
	db.ID(uauthor).Update(&user)
}

func ChangeName(c *gin.Context) {
	uauthor := GetUserIdWithoutCheck(c)
	name := c.PostForm("name")
	var user User
	db.ID(uauthor).Get(&user)
	user.Name = name
	user.Level = user.Level - 8
	db.ID(uauthor).Update(&user)
}

func VerifyKey(key string) bool {
	if len(key) != 6 {
		return false
	}
	//判断第一位是否为数字
	sN := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "00"}
	for _, i := range sN {
		if i == "00" {
			return false
		}
		if strings.HasPrefix(key, i) {
			break
		}
	}
	exist, _ := db.In("key", key).Exist(&KeyUsed{})
	if exist {
		return false
	}
	d := []byte(key[1:])
	m := md5.New()
	m.Write(d)
	assign := hex.EncodeToString(m.Sum(nil))
	if strings.Contains(assign, assginSegment) {
		return true
	} else {
		return false
	}
}
