package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	store := memstore.NewStore([]byte("just_secret"))
	r.Use(sessions.Sessions("session_id", store))
	r.GET("/ping", func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("count")
		var count int
		if v == nil {
			count = 0
		} else{
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"message": count,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
