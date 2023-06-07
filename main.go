package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-contrib/sessions"
	sred "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
	
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

	store, err := sred.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		fmt.Println("redis坏掉了😵")
		panic(err)
	}
	store.Options(sessions.Options{
		Secure:   true, //跟下面那条基本上可以防住csrf了,但是还是稳一点好
		HttpOnly: true,
		Path:     "/",
		MaxAge:   1200000}) //凑个整,差一点点到2week
	r.Use(sessions.Sessions("session_id", store))
	r.LoadHTMLGlob("templates/**/*")
	// TODO csrf防护,需要前端支持

	
	db, err = xorm.NewEngine("postgres", "postgresql://postgres:114514@localhost:5432/dbs?sslmode=disable")
	if err != nil {
		panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
	}

	db.Sync(&User{}, &Video{}, &VideoComment{}, &Tag{}, &Forum{})
	db.SetDefaultCacher(caches.NewLRUCacher(caches.NewMemoryStore(), 1000))


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
		group.POST("/main_forum", func(ctx *gin.Context) {
			strFid := ctx.PostForm("fid")
			intFid, _ := strconv.Atoi(strFid)
			uintFid := uint(intFid)
			FindMainForum(uintFid, ctx)
		})
		group.POST("/unit_forum", func(ctx *gin.Context) {
			strFid := ctx.PostForm("fid")
			intFid, _ := strconv.Atoi(strFid)
			uintFid := uint(intFid)
			strId := ctx.PostForm("id")
			intId, _ := strconv.Atoi(strId)
			uintId := uint(intId)
			FindUnitForum(uintFid, uintId, ctx)
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
