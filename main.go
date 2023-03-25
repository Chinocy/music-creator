package main

import (
	"music-creator/internal/webserver"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	webserver.Run()
}
