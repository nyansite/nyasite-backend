// remember to creat token.go !!!!!
package main

import (
	"context"
	"crypto/rand"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	db.Sync(&User{}, &Tag{}, &TagModel{}, &SessionSecret{},
		&VideoNeedToCheck{}, &Video{}, &VideoPlayedRecord{},
		&VideoComment{}, &VideoCommentReply{}, &VideoCommentEmojiRecord{}, &VideoCommentReplyLikeRecord{},
		&VideoBullet{},
		&Invitation{}, &Discharge{})
	db.SetDefaultCacher(caches.NewLRUCacher(caches.NewMemoryStore(), 100000))
	//上面的是sql

	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://google.com"}	//允许访问信息的第三方,比如说广告供应商
	// config.AllowCredentials = true //cookie一并发给跨域请求
	// r.Use(cors.New(config))
	var secrets [][]byte
	var old_secrets []SessionSecret
	s1, _ := rand.Prime(rand.Reader, 256) //最多32字节,也就是256
	s2, _ := rand.Prime(rand.Reader, 256)
	secrets = append(secrets, s1.Bytes(), s2.Bytes())

	//大概需要平均每个TTL周期重启一次,因为我懒得拉一个goroutine //TODO 自动更新密钥
	db.Where("created_at < ?", time.Now().Unix()-TTL).Delete(&SessionSecret{})  //删除过期
	err = db.Where("created_at >= ?", time.Now().Unix()-TTL).Find(&old_secrets) //没过期的取出来
	if err != nil {
		panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!") //数据库连不上会在这里挂,而不是上面
	}
	db.Insert(&SessionSecret{Authentication: s1.Bytes(), Encryption: s2.Bytes()}) //新密钥进数据库,不用defer避免kill 9
	for _, v := range old_secrets {
		secrets = append(secrets, v.Authentication, v.Encryption)
	}
	store := cookie.NewStore(secrets...) //密钥成对定义以允许密钥轮换.使用新的密钥加密但是旧的仍然有效
	store.Options(sessions.Options{
		Secure:   true, //跟下面那条基本上可以防住csrf了,但是还是稳一点好 //TODO 更靠谱的csrf防护
		HttpOnly: true, //localhost或者https
		Path:     "/",
		MaxAge:   TTL,
		SameSite: http.SameSiteStrictMode})

	//下面是路由
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(defaultLogFormatter), gin.Recovery())
	r.Use(sessions.Sessions("session", store))
	r.MaxMultipartMemory = 8 << 20 //8mb,默认32,限制每个请求的内存占用,但是不会影响接收大文件
	group := r.Group("/api")
	{
		group.GET("/coffee", CheckPrivilege(11), coffee)
		group.GET("/user_status", GetSelfUserData)
		group.GET("/user_status/:id", GetUserData)

		group.POST("/register", Register)
		group.POST("/login", Login)

		group.POST("/upload_img", LimitRequestBody(2<<20), PostImg)

		group.GET("/taglist", EnireTag)
		group.POST("/new_tag", CheckPrivilege(3), NewTag)
		//video
		group.GET("/get_video/:id", GetVideo)
		group.POST("/upload_video", CheckPrivilege(1), PostVideo)
		group.GET("/get_video_tags/:id", GetVideoTags)
		//comment
		group.GET("/video_comment/:id/:pg", BrowseVideoComments)
		group.GET("/video_comment_reply/:id", BrowseVideoCommentReplies)
		group.POST("/add_video_comment", CheckPrivilege(1), AddVideoComment)
		group.POST("/add_video_comment_reply", CheckPrivilege(1), AddVideoCommentReply)
		group.POST("/click_comment_emoji", ClikckCommentEmoji)
		group.POST("/click_commentreply_like", ClickCommentReplyLike)
		group.POST("/add_video_tag", CheckPrivilege(10), AddVideoTag)
		//danmaku
		group.GET("/get_bullets/:id", BrowseBullets)
		group.POST("/add_video_bullet", CheckPrivilege(0), AddBullet)
		//change user information
		group.POST("/change_avatar", CheckPrivilege(0), ChangeAvatar)
		group.POST("/change_name", CheckPrivilege(1), ChangeName)
		//search
		group.POST("/search_video", SearchVideos)

		group.GET("/get_all_videos_needtocheck", CheckPrivilege(10), GetAllVideoNeedToChenck)
		group.POST("/pass_video", CheckPrivilege(10), PassVideo)
		group.POST("/reject_video", CheckPrivilege(10), RejectVideo)

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

	quit := make(chan os.Signal, 1)
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
