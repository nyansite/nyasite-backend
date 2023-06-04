package main

import (
	"fmt"

	"github.com/gin-contrib/cors"

	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/redis/go-redis/v9"

	// "github.com/gin-contrib/sessions/memstore"
	"time"

	sred "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var(
	db *gorm.DB
	Tags []string
	rdb *redis.Client
)


func main() {
	r := gin.Default()
	// store := cookie.NewStore([]byte("just_secret")) //不安全但是方便测试,记得清cookie
	store, err := sred.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		fmt.Println("redis坏掉了😵")
		panic(err)
	}
	store.Options(sessions.Options{Secure: true, HttpOnly: true, Path: "/", MaxAge: 3000000})
	r.Use(sessions.Sessions("session_id", store))
	r.LoadHTMLGlob("templates/**/*")
	// TODO csrf防护,需要前端支持

	db, err = gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{
		PrepareStmt: true, //执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
	})
	if err != nil {
		panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
	}
	db.AutoMigrate(&User{}, &Video{}, &VideoComment{}, &Tag{}, &Forum{}, &ForumComment{}) //实际上的作用是创建表

	rdb = redis.NewClient(&redis.Options{
		Addr:	  "localhost:6379",
		Password: "", // no password set
		DB:		  0,  // use default DB
	})

	group := r.Group("/api")
	{
		group.GET("/user_status", GetSelfUserData)
		group.GET("/user_status/:id", GetUserData)
		group.GET("/video_comment/:id/:pg", GetVideoComment)
		group.GET("/video_img/:id", GetVideoImg)
		group.GET("/coffee", coffee)

		group.POST("/register", Register)
		group.POST("/login", Login)
		group.POST("/new_tag", NewTag)
		group.POST("/add_comment", AddComment)
		group.POST("/upload_video", UploadVideo)
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
		group.Static("img", "./img")
	}
	//管理员页面
	group = r.Group("/admin")
	group.Use(AdminCheck())
	{
		group.GET("/browse_video/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "browsevideo.html", gin.H{})
		})
		group.GET("/upload_video", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "uploadvideo.html", gin.H{})
		})

		group.POST("/browse_video/:page", AdminVideoPost)
		group.POST("/upload_video", UploadVideo)
	}

	// db.Create(&Video{})
	// var i uint64
	// for i = 0; i < 114; i++ {
	// 	db.Create(&Video{})
	// }
	// rdb.Set(context.Background(), "1", 100, 0)
	// val, err := rdb.Get(context.Background(), "1").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("返回", val)
	r.Run(":8000") // 8000
}

func coffee(c *gin.Context) { //没有人能拒绝愚人节彩蛋
	if time.Now().Month() == 4 && time.Now().Day() == 1 {
		c.String(http.StatusTeapot, "我拒绝泡咖啡,因为我是茶壶")
	} else {
		c.String(http.StatusForbidden, "我拒绝泡咖啡,因为我是服务器")
	}
}
