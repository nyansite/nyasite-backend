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
	store.Options(sessions.Options{Secure: true, HttpOnly: true})
	r.Use(sessions.Sessions("session_id", store))
	// TODO csrf防护,需要前端支持

	// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// if err != nil {
	// 	panic("我数据库呢???我那么大一个数据库呢???还我数据库!!!")
	// }

	r.GET("/api/user_status", get_self_user_status)
	r.GET("/api/user_status/:id", get_self_user_status)


	r.POST("/api/register", post_register)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}




func get_self_user_status(c *gin.Context) {
	
	session := sessions.Default(c)
	userid := 0 //userid为零表示错误
	if session.Get("is_login") == true {
		userid = session.Get("userid").(int)
	// } else {
	// 	c.AbortWithStatus(http.StatusUnauthorized) //返回401
	// 	return
	}

	c.JSON(http.StatusOK, gin.H{
		"userid": userid,
	})
}

func get_user_status(c *gin.Context)  {
	// TODO 根据id获取
}

func login(c *gin.Context) {

}


func post_register(c *gin.Context) {
	
	session := sessions.Default(c)
	csrf_token := session.Get("csrf_token")
	if csrf_token == nil{
		return
	}
	c.String(http.StatusOK, csrf_token.(string))
}
