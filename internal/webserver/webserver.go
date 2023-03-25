package webserver

import (
	"music-creator/internal/webserver/song"

	"github.com/gin-gonic/gin"
)

type config struct {
	Token string
}

func Run() {
	r := gin.Default()
	songCtrl := song.NewSongCtrl()
	r.POST("/song", songCtrl.CreateSong)
	r.GET("/song/choices", songCtrl.GetChoices)
	r.Run()
}
