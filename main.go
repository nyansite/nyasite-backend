package main

import (
	// "cute_site/models"

	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
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

func login() {

}

func register() {

}
