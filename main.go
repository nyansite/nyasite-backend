package main

import (
	"cute_site/models"

	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"crypto/rand"
	"crypto/sha512"
	"regexp"
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
	db.AutoMigrate(&models.User{}) //实际上的作用是创建表

	group := r.Group("/api")
	{
		group.GET("/user_status", get_self_user_status)
		group.GET("/user_status/:id", get_user_status)

		group.POST("/register", register)
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func get_self_user_status(c *gin.Context) {

	session := sessions.Default(c)
	userid := 0 //userid为零表示错误
	if session.Get("is_login") == true {
		userid = session.Get("userid").(int)
		// } else {
		// 	c.AbortWithStatus(http.StatusUnauthorized) //返回401
		// 	return
	}

	c.JSON(http.StatusOK, gin.H{
		"userid": userid,
	})
}

func get_user_status(c *gin.Context) {
	// TODO 根据id获取
}

func login(c *gin.Context) {

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

	user := models.User{Name: username, Passwd: encrypt_passwd([]byte(username)), Email: mail}
	db.Create(&user)
	c.AbortWithStatus(601)
}

func encrypt_passwd(passwd []byte) []byte { //加密密码,带盐
	salte, _ := rand.Prime(rand.Reader, 64) //普普通通的64位盐
	salt := salte.Bytes()

	passwd_sha := sha512.Sum512(passwd)          //密码的sha
	saltpasswd := append(passwd_sha[:], salt...) //加盐
	safe_passwd := sha512.Sum512(saltpasswd)     //这一步才算加密

	return append(safe_passwd[:], salt...) //保存盐
}

func check_passwd(passwd []byte, passwd2 []byte) bool {
	//获取盐

	return true
}
