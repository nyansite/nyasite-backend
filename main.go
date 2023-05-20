package main

import (
	"cute_site/models"
	"fmt"

	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	_ "encoding/json"
	_ "fmt"
	"regexp"
	"time"
)

var db *gorm.DB

func main() {
	r := gin.Default()
	store := memstore.NewStore([]byte("just_secret"))
	store.Options(sessions.Options{Secure: true, HttpOnly: true})
	r.Use(sessions.Sessions("session_id", store))
	// TODO csrf防护,需要前端支持

	db_l, dberr := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db = db_l
	if dberr != nil {
		panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
	}
	
	db.AutoMigrate(&models.User{}, &models.Video{}, &models.Comment{}) //实际上的作用是创建表
	
	group := r.Group("/api")
	{
		group.GET("/user_status", get_self_user_status)
		group.GET("/user_status/:id", get_user_status)
		group.GET("/coffee", coffee)

		group.POST("/register", register)
		group.POST("/login", login)
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func get_self_user_status(c *gin.Context) {

	session := sessions.Default(c)
	var user models.User

	if session.Get("is_login") == true {
		userid := session.Get("userid")
		level := session.Get("level")
		db.First(&user, userid)
		mail := user.Email
		c.JSON(http.StatusOK, gin.H{
			"userid": userid,
			"mail":   mail,
			"level": level,
		})
	} else {
		c.AbortWithStatus(http.StatusUnauthorized) //返回401
		return
	}
}

func get_user_status(c *gin.Context) {
	// TODO 根据id获取
}

func coffee(c *gin.Context) { //没有人能拒绝愚人节彩蛋
	if time.Now().Month() == 4 && time.Now().Day() == 1 {
		c.String(http.StatusTeapot, "我拒绝泡咖啡,因为我是茶壶")
	} else {
		c.String(http.StatusForbidden, "我拒绝泡咖啡,因为我是服务器")
	}
}

//下面的都是post

func login(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("is_login") == true {
		c.AbortWithStatus(610)
		return
	}
	username, passwd := c.PostForm("username"), c.PostForm("passwd") //传入的用户名也有可能是邮箱
	if username == "" || passwd == "" {
		c.AbortWithStatus(http.StatusBadRequest) //400
		return
	}
	var user models.User
	if db.First(&user, "Name = ? OR Email = ?", username, username).RowsAffected == 0 { //用户不存在
		c.AbortWithStatus(612)
		return
	}
	if !check_passwd(user.Passwd, []byte(passwd)) {
		c.AbortWithStatus(613)
		return
	}

	session.Set("userid", user.ID)
	session.Set("is_login", true)
	session.Set("level", user.Level)
	session.Save()
	c.AbortWithStatus(611)
}

func register(c *gin.Context) {
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

	//601:创建成功,602:用户名重复,603:邮箱重复
	if db.First(&models.User{}, "Name = ?", username).RowsAffected != 0 {
		c.AbortWithStatus(602)
		return
	}
	if db.First(&models.User{}, "Email = ?", mail).RowsAffected != 0 {
		c.AbortWithStatus(603)
		return
	}

	user := models.User{Name: username, Passwd: encrypt_passwd([]byte(passwd)), Email: mail}
	db.Create(&user)
	c.AbortWithStatus(601)
}

//上面的都是post

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
