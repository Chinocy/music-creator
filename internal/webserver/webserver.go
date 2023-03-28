package webserver

import (
	"music-creator/internal/webserver/song"

	"github.com/gin-gonic/gin"
)

type config struct {
	Token string
}

func Run() {
	songCtrl := song.NewSongCtrl()

	r := gin.Default()
	r.POST("/song", songCtrl.CreateSong)
	r.GET("/song/choices", songCtrl.GetChoices)
	r.Run()
}
