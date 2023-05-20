package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "encoding/json"
	_ "fmt"
	"time"
)

var db *gorm.DB
var tags []TagText

func main() {
	r := gin.Default()
	store := memstore.NewStore([]byte("just_secret"))
	store.Options(sessions.Options{Secure: true, HttpOnly: true})
	r.Use(sessions.Sessions("session_id", store))
	// TODO csrf防护,需要前端支持

	tags = []TagText{}
	dbl, dberr := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db = dbl
	if dberr != nil {
		panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
	}

	db.AutoMigrate(&User{}, &Video{}, &Comment{}, &Tag{}) //实际上的作用是创建表

	group := r.Group("/api")
	{
		group.GET("/user_status", GetSelfUserData)
		group.GET("/user_status/:id", GetUserData)
		group.GET("/coffee", coffee)

		group.POST("/register", Register)
		group.POST("/login", Login)

	}

	r.Run() // 8080

}

func coffee(c *gin.Context) { //没有人能拒绝愚人节彩蛋
	if time.Now().Month() == 4 && time.Now().Day() == 1 {
		c.String(http.StatusTeapot, "我拒绝泡咖啡,因为我是茶壶")
	} else {
		c.String(http.StatusForbidden, "我拒绝泡咖啡,因为我是服务器")
	}
}
