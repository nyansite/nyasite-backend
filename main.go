package main

import (
	// "cute_site/models"

	
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"

	"crypto/rand"
	"encoding/base64"
)

func main() {
	r := gin.Default()
	store := memstore.NewStore([]byte("just_secret"))
	r.Use(sessions.Sessions("session_id", store))
	// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// if err != nil {
	// 	panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
	// }

	r.GET("/api/user_status", get_user_status)
	r.POST("/api/register", post_register)

	r.GET("/uapi/csrf_token", get_csrf_token) //不安全的api,记得nginx禁用
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func get_user_status(c *gin.Context) {
	session := sessions.Default(c)
	userid := 0 //userid为零表示错误
	if session.Get("is_login") == true {
		userid = session.Get("userid").(int)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized) //返回401
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userid": userid,
	})
}

func login(c *gin.Context) {

}

func get_csrf_token(c *gin.Context) { //获取csrf token,被表单携带
	session := sessions.Default(c)
	n, _ := rand.Prime(rand.Reader, 128)

	csrf_token := base64.StdEncoding.EncodeToString(n.Bytes())
	session.Set("csrf_token", csrf_token)
	session.Save()
	c.String(http.StatusOK, csrf_token)
}

func post_register(c *gin.Context) {
	session := sessions.Default(c)
	csrf_token := session.Get("csrf_token")
	if csrf_token == nil{
		return
	}
	c.String(http.StatusOK, csrf_token.(string))
}
