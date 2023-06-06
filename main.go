package main

import (
	"github.com/gin-contrib/sessions"
	"net/http"
	"time"
	// "github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	xcach "xorm.io/xorm/caches"
	"xorm.io/xorm/log"
)

var (
	db *xorm.Engine
)

func main() {
	r := gin.Default()
	
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://google.com"}	//允许访问信息的第三方,比如说广告供应商
	// config.AllowCredentials = true	//cookie一并发给跨域请求
	// r.Use(cors.New(config))

	store := cookie.NewStore([]byte("just_secret")) //不安全但是方便测试,记得清cookie
	// store := memstore.NewStore([]byte("secret"))

	store.Options(sessions.Options{
		Secure:   true, //跟下面那条基本上可以防住csrf了,但是还是稳一点好
		HttpOnly: true,
		Path:     "/",
		MaxAge:   1000000}) //大概不到12d
	r.Use(sessions.Sessions("session_id", store))
	r.LoadHTMLGlob("templates/**/*")
	// TODO csrf防护,需要前端支持

	var err error
	db, err = xorm.NewEngine("sqlite3", "./test.db")
	if err != nil {
		panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
	}
	db.Logger().SetLevel(log.LOG_INFO)
	db.Sync(&User{}, &Video{}, &VideoComment{}, &Tag{}, &Forum{}, &ForumComment{})
	db.SetDefaultCacher(xcach.NewLRUCacher(xcach.NewMemoryStore(), 1000))

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

	r.Run(":8000") // 8000
}

func coffee(c *gin.Context) { //没有人能拒绝愚人节彩蛋
	if time.Now().Month() == 4 && time.Now().Day() == 1 {
		c.String(http.StatusTeapot, "我拒绝泡咖啡,因为我是茶壶")
	} else {
		c.String(http.StatusForbidden, "我拒绝泡咖啡,因为我是服务器")
	}
}
