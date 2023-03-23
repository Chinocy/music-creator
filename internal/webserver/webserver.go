package webserver

import (
	"log"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
	songCtrl := song.NewSongCtrl()
	r.POST("/song", songCtrl.CreateSong)
	r.Run()
}
