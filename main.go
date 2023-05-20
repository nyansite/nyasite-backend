package main

import (
	"cute_site/models"
	"cute_site/user"
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

func main() {
	r := gin.Default()
	store := memstore.NewStore([]byte("just_secret"))
	store.Options(sessions.Options{Secure: true, HttpOnly: true})
	r.Use(sessions.Sessions("session_id", store))
	// TODO csrf防护,需要前端支持

	db, dberr := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if dberr != nil {
		panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
	}

	db.AutoMigrate(&models.User{}, &models.Video{}, &models.Comment{}, &models.Tag{}) //实际上的作用是创建表
	// tags := []models.TagText{}
	group := r.Group("/api")
	{
		group.GET("/user_status", func(ctx *gin.Context) { user.GetSelfUserStatus(ctx, db) })
		group.GET("/user_status/:id", func(ctx *gin.Context) { user.GetUserStatus(ctx, db) })
		group.GET("/coffee", coffee)

		group.POST("/register", func(ctx *gin.Context) { user.Register(ctx, db) })
		group.POST("/login", func(ctx *gin.Context) { user.Login(ctx, db) })

	}

	r.Run(":8000") // listen and serve on 0.0.0.0:8000 (for windows "localhost:8000")

}

func coffee(c *gin.Context) { //没有人能拒绝愚人节彩蛋
	if time.Now().Month() == 4 && time.Now().Day() == 1 {
		c.String(http.StatusTeapot, "我拒绝泡咖啡,因为我是茶壶")
	} else {
		c.String(http.StatusForbidden, "我拒绝泡咖啡,因为我是服务器")
	}
}
