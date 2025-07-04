// remember to creat token.go !!!!!
package main

import (
	"context"
	"crypto/rand"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"google.golang.org/api/gmail/v1"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
)

var (
	db       *xorm.Engine
	emailSrv *gmail.Service
)

const (
	TTL = 1200000 //session的寿命,单位秒,接近两周
)

func main() {
	var err error
	db, _ = xorm.NewEngine("postgres", "postgresql://mbe:"+databasePassword+"@localhost:5432/dbs?sslmode=disable")
	err = db.Ping()
	if err != nil {
		log.Println("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
		panic(err)
	}
	c := cron.New()
	//创建定时任务对象（热榜）
	db.Sync(&User{}, &Tag{}, &TagModel{}, &SessionSecret{},
		&VideoNeedToCheck{}, &Video{}, &VideoPlayedRecord{}, &VideoLikeRecord{}, &VideoMarkRecord{},
		&VideoComment{}, &VideoCommentReply{}, &VideoCommentEmojiRecord{}, &VideoCommentReplyLikeRecord{},
		&VideoBullet{},
		&Circle{}, &MemberOfCircle{}, &ApplyCircle{}, &VoteOfApplyCircle{},
		&Invitation{}, &Discharge{},
		&TrendingRankVideo{},
		&KeyUsed{})
	db.SetDefaultCacher(caches.NewLRUCacher(caches.NewMemoryStore(), 10000))
	//上面的是sql
	VerCodeAllocMap = make(map[string]VerCode)
	//上面的是初始化验证码队列
	emailSrv, err = GenSrv()
	if err != nil {
		panic(err)
	}
	//
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://google.com"}	//允许访问信息的第三方,比如说广告供应商
	// config.AllowCredentials = true //cookie一并发给跨域请求
	// r.Use(cors.New(config))
	var secrets [][]byte
	var old_secrets []SessionSecret
	s1, _ := rand.Prime(rand.Reader, 256) //最多32字节,也就是256
	s2, _ := rand.Prime(rand.Reader, 256)
	secrets = append(secrets, s1.Bytes(), s2.Bytes())
	db.Where("created_at < ?", time.Now().Unix()-TTL).Delete(&SessionSecret{})  //删除过期
	err = db.Where("created_at >= ?", time.Now().Unix()-TTL).Find(&old_secrets) //没过期的取出来
	if err != nil {
		panic(err) //很难相信这里会出问题
	}
	db.Insert(&SessionSecret{Authentication: s1.Bytes(), Encryption: s2.Bytes()}) //新密钥进数据库,不用defer避免kill 9
	for _, v := range old_secrets {
		secrets = append(secrets, v.Authentication, v.Encryption)
	}
	store := cookie.NewStore(secrets...) //密钥成对定义以允许密钥轮换.使用新的密钥加密但是旧的仍然有效
	store.Options(sessions.Options{
		Secure:   true, //跟下面那条基本上可以防住csrf了,但是还是稳一点好// TODO 更靠谱的csrf防护
		HttpOnly: true, //localhost或者https
		Path:     "/",
		MaxAge:   TTL,
		SameSite: http.SameSiteStrictMode})
	//下面是定时任务
	c.AddFunc("@daily", RankVideosDaliy)
	c.AddFunc("@mouthly", RankVideosMouthly)
	c.AddFunc("@yearly", RankVideosYearly)
	//下面是路由
	r := gin.Default()
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		log.Fatal(err)
	}
	r.Use(gin.Recovery(), sessions.Sessions("session", store))
	group := r.Group("/api")
	{
		group.GET("/user_status", GetSelfUserData)
		group.GET("/user_status/:id", GetUserData)

		group.GET("/search_users/:name", CheckPrivilege(0), SearchUsers)

		group.POST("/register", Register)
		group.GET("/logout", QuitLogin)
		group.POST("/login", Login)
		group.POST("/reset_pwd", ResetPwd)
		group.GET("/refresh", Refresh)

		group.POST("/clockin", ClockIn)

		group.GET("/check_premission/:cid", CheckPrivilege(0), CheckPremissionOfCircle)

		group.GET("/coffee", CheckPrivilege(11), coffee)
		group.GET("/taglist", EnireTag)

		group.POST("/get_ver_code_reset_pwd", AllocVerCodeResetPwd)
		group.POST("/get_ver_code_register", AllocVerCodeRegister)
		//trendingTest
		group.POST("/trendingTest", RankVideosTest)

		group.POST("/new_tag", CheckPrivilege(10), NewTag)
		//video
		group.POST("/get_upload_url", UploadVideoUrl)
		group.GET("/get_video/:id", GetVideo)
		group.POST("/upload_video", CheckPrivilege(0), PostVideo)
		group.GET("/get_all_videos", GetAllVideos)
		group.GET("/get_video_tags/:id", GetVideoTags)
		group.GET("/get_video_link/:uid", GetLinkPlaybackHls)
		group.POST("/like_video", CheckPrivilege(0), LikeVideo)
		group.POST("/mark_video", CheckPrivilege(0), MarkVideo)
		//comment
		group.GET("/video_comment/:id/:pg", BrowseVideoComments)
		group.GET("/video_comment_reply/:id", BrowseVideoCommentReplies)
		group.POST("/add_video_comment", CheckPrivilege(0), AddVideoComment)
		group.POST("/add_video_comment_reply", CheckPrivilege(0), AddVideoCommentReply)
		group.POST("/click_comment_emoji", ClikckCommentEmoji)
		group.POST("/click_commentreply_like", ClickCommentReplyLike)
		group.POST("/add_video_tag", CheckPrivilege(9), AddVideoTag)
		//danmaku
		group.GET("/get_bullets/:id", BrowseBullets)
		group.POST("/add_video_bullet", CheckPrivilege(0), AddBullet)
		//change user information
		group.POST("/change_avatar", CheckPrivilege(0), ChangeAvatar)
		group.POST("/change_name", CheckPrivilege(1), ChangeName)
		//circle
		group.POST("/apply_circle", CheckPrivilege(0), PostCircleApplication)
		group.GET("/get_available_circle/:type", CheckPrivilege(0), CheckAvailableCircle)
		group.GET("/get_circle/:id", GetCircle)
		group.GET("/get_circle_joined", CheckPrivilege(0), GetCircleJoined)
		group.GET("/get_circle_subscribed", CheckPrivilege(0), GetCirclesSubscribed)
		group.POST("/subscribe", CheckPrivilege(0), SubscribeCircle)
		group.GET("/get_circle_video/:id/:page/:method", GetVideosOfCircle)
		//circle manage
		group.GET("/get_circle_members/:cid", CheckPrivilege(0), GetAllMembersOfCircle)
		group.POST("/invite_new_member", CheckPrivilege(0), InviteMember)
		group.POST("/kick_out", CheckPrivilege(0), KickOut)
		group.POST("/delete_video", CheckPrivilege(0), DeleteVideo)
		//search
		group.POST("/search_video", SearchVideos)
		//token
		group.GET("/get_HELLOIMG_token", CheckPrivilege(0), GetHELLOIMGtoken)
		//check
		group.GET("/get_all_circles_needtocheck", CheckPrivilege(9), GetAllCirclesNeedtoCheck)
		group.POST("/vote_for_circles_needtocheck", CheckPrivilege(9), VoteForCirclesNeedtoCheck)
		group.GET("/get_all_videos_needtocheck", CheckPrivilege(9), GetAllVideoNeedToChenck)
		group.POST("/pass_video", CheckPrivilege(9), PassVideo)
		group.POST("/reject_video", CheckPrivilege(9), RejectVideo)
		group.POST("/withdraw_video", CheckPrivilege(9), WithdrawVideo)
		//message
		group.GET("/get_video_subscribed", CheckPrivilege(0), GetVideoSubscribe)
		group.GET("/get_new_circle_affairs", CheckPrivilege(0), GetCircleAffairs)
		group.POST("/reply_invitation", CheckPrivilege(0), ReplyInvitation)
		group.GET("/get_check_messages", CheckPrivilege(0), GetCheckMessage)
		group.POST("/delete_video_need_to_check", CheckPrivilege(0), DeleteVideoNeedToCheck)
		//trending
		group.GET("/get_trending", GetTrending)
		//history
		group.GET("/history/:pg", CheckPrivilege(0), GetHistoryRecord)
		//marks
		group.GET("/marks/:pg", CheckPrivilege(0), GetMarksRecord)
	}

	//  https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	srv := &http.Server{
		Addr:    ":8000",
		Handler: r.Handler(),
	}
	go func() {
		log.Println("服务器启动")
		c.Start()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL 无法被捕获,也不需要捕获
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit //等待信号,阻塞

	log.Println("服务器关闭中~~~")
	c.Stop()
	ctx, channel := context.WithTimeout(context.Background(), 5*time.Second) //关闭等几秒,多线程很有用
	defer channel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器关闭错误(不用管):", err)
	}

	// <-ctx.Done()//不想等就把这行干掉
	log.Println("服务器关闭")
}

func coffee(c *gin.Context) { //没有人能拒绝愚人节彩蛋
	if time.Now().Month() == 4 && time.Now().Day() == 1 {
		c.String(http.StatusTeapot, "我拒绝泡咖啡,因为我是茶壶")
	} else {
		c.String(http.StatusServiceUnavailable, "我拒绝泡咖啡,因为我是服务器")
	}
}
