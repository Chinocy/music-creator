package webserver

import (
	"music-creator/internal/webserver/song"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type config struct {
	Token string
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header["Authorization"]
		if token[0] != os.Getenv("AUTH_TOKEN") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Run() {
	songCtrl := song.NewSongCtrl()

	r := gin.Default()
	r.Use(Authorize())
	r.POST("/song", songCtrl.CreateSong)
	r.GET("/song/choices", songCtrl.GetChoices)
	r.Run()

	s := &http.Server{
		Handler:        r,
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
