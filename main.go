package main

import (
	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/cookie"
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
	// store := memstore.NewStore([]byte("just_secret"))
	store := cookie.NewStore([]byte("just_secret")) //不安全但是方便测试,记得清cookie
	store.Options(sessions.Options{Secure: true, HttpOnly: true})
	r.Use(sessions.Sessions("session_id", store))
	r.LoadHTMLGlob("test/*")
	// TODO csrf防护,需要前端支持

	tags = []TagText{}
	dbl, dberr := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db = dbl
	if dberr != nil {
		panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
	}

	db.AutoMigrate(&User{}, &Video{}, &VideoCommentPage{}, &VideoComment{},&VideoCommentSubPage{},&VideoCommentReply{}, &Tag{}, &TagText{}) //实际上的作用是创建表

	group := r.Group("/api")
	{
		group.GET("/user_status", GetSelfUserData)
		group.GET("/user_status/:id", GetUserData)
		group.GET("/video_comment/:id/:pg", GetVideoComment)
		group.GET("/coffee", coffee)

		group.POST("/register", Register)
		group.POST("/login", Login)
		group.POST("/new_tag", NewTag)
		group.POST("/add_comment", AddComment)
		
	}
	config := cors.Config{
		AllowOrigins: []string{"https://127.0.0.1"}, //只允许本地访问
	} //这个是不允许远程的
	group = r.Group("/uapi") //不安全的api,能够操作数据库的所有数据
	group.Use(cors.New(config))
	{

	}

	group = r.Group("/test")
	{
		group.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{})
		})
		group.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.html", gin.H{})
		})
		group.GET("/register", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "register.html", gin.H{})
		})
		group.GET("/add_file", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "addfile.html", gin.H{})
		})
		group.GET("get_file", GetFile)
		group.POST("/api/add_file", AddFile)
		group.Static("img", "./img")
	}

	r.Run(":8000") // 8000
	// db.Create(&Video{CommentP: []CommentPage{{Comment: []Comment{{Text: "ww"}}}}})
	// var i uint64
	// for i = 0; i < 10; i++ {
	// 	WAddComment("只因", 1, 0)
	// }
}

func coffee(c *gin.Context) { //没有人能拒绝愚人节彩蛋
	if time.Now().Month() == 4 && time.Now().Day() == 1 {
		c.String(http.StatusTeapot, "我拒绝泡咖啡,因为我是茶壶")
	} else {
		c.String(http.StatusForbidden, "我拒绝泡咖啡,因为我是服务器")
	}
}
