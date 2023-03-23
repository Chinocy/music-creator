package webserver

import (
	"music-creator/internal/webserver/song"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type config struct {
	Token string
}

func Run() {
	// gin.SetMode(gin.DebugMode)
	r := gin.Default()
	godotenv.Load()
	songCtrl := song.NewSongCtrl()
	r.POST("/song", songCtrl.CreateSong)
	r.GET("/song/choices", songCtrl.GetChoices)
	r.Run()
}
