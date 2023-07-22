package main

import (
	"context"
	"crypto/rand"
	"cute_site/static"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
)

var (
	db *xorm.Engine
)

const (
	TTL = 1200000 //session的寿命,单位秒,接近两周
)

func main() {
	var err error
	db, err = xorm.NewEngine("postgres", "postgresql://postgres:114514@localhost:5432/dbs?sslmode=disable")
	if err != nil {
		panic(err) //连接失败不会在这里挂
	}

	db.Sync(&User{}, &Video{}, &VideoComment{}, &Tag{}, &Forum{}, &SessionSecret{}, &ForumComment{}, &EmojiRecord{})
	db.SetDefaultCacher(caches.NewLRUCacher(caches.NewMemoryStore(), 1000))
	//上面的是sql

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} //根据个人的实例来更改端口
	config.AllowHeaders = []string{"Origin", "X-Requested-With", "Content-Type",
		"Accept", "Authorization", "Access-Control-Allow-Origin"}
	config.AllowCredentials = true //cookie一并发给跨域请求
	r.Use(cors.New(config))
	var secrets [][]byte
	var old_secrets []SessionSecret
	s1, _ := rand.Prime(rand.Reader, 256) //最多32字节,也就是256
	s2, _ := rand.Prime(rand.Reader, 256)
	secrets = append(secrets, s1.Bytes(), s2.Bytes())
	db.Where("created_at < ?", time.Now().Unix()-TTL).Delete(&SessionSecret{}) //删除过期
	err = db.Where("created_at >= ?", time.Now().Unix()-TTL).Find(&old_secrets)
	if err != nil {
		panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
	}
	db.Insert(&SessionSecret{Authentication: s1.Bytes(), Encryption: s2.Bytes()})
	for _, v := range old_secrets {
		secrets = append(secrets, v.Authentication, v.Encryption)
	}
	store := cookie.NewStore(secrets...)
	store.Options(sessions.Options{
		//Secure:   true, //跟下面那条基本上可以防住csrf了,但是还是稳一点好
		HttpOnly: true, //测试阶段调ssl有点麻烦
		Path:     "/",
		MaxAge:   TTL,
		SameSite: http.SameSiteStrictMode})
	r.Use(sessions.Sessions("session", store))
	r.Use(static.Serve("/", static.LocalFile("cute_web/build/", false)))
	// TODO csrf防护,需要前端支持
	group := r.Group("/api")
	{
		group.GET("/user_status", GetSelfUserData)
		group.GET("/user_status/:id", GetUserData)
		group.GET("/get_video_img/:id", GetVideoImg)
		group.GET("/video_comment/:id/:pg", GetVideoComment)
		group.GET("/coffee", PrivilegeLevel(11), coffee)
		group.GET("/all_forum/:page", BrowseAllForumPost)
		group.GET("/browse_forum/:board/:page", BrowseForumPost)
		group.GET("/browse_unitforum/:mid/:page", BrowseUnitforumPost)

		group.POST("/register", Register)
		group.POST("/login", Login)
	}

	group = r.Group("/uapi")
	{
		group.POST("/new_tag", PrivilegeLevel(10), NewTag)
		//video
		group.POST("/upload_video", UploadVideo)
		group.POST("/add_video_comment", AddVideoComment)
		//fourm
		group.POST("/add_mainforum", PrivilegeLevel(0), AddMainforum)
		group.POST("/add_unitforum", PrivilegeLevel(0), AddUnitforum)
		group.POST("/add_emoji", PrivilegeLevel(0), AddEmoji)
		group.POST("/finish_forum", PrivilegeLevel(0), FinishForum)
	}

	//  https://gin-gonic.com/zh-cn/docs/examples/graceful-restart-or-stop/
	srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}
	go func() {
		log.Println("服务器启动")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit //等待信号,阻塞

	log.Println("服务器关闭中~~~")

	ctx, channel := context.WithTimeout(context.Background(), 5*time.Second)
	defer channel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器关闭错误(不用管):", err)
	}
	log.Println("服务器关闭")
}

func coffee(c *gin.Context) { //没有人能拒绝愚人节彩蛋
	if time.Now().Month() == 4 && time.Now().Day() == 1 {
		c.String(http.StatusTeapot, "我拒绝泡咖啡,因为我是茶壶")
	} else {
		c.String(http.StatusServiceUnavailable, "我拒绝泡咖啡,因为我是服务器")
	}
}
