package test

import "github.com/joho/godotenv"

func SetupTest() {
	godotenv.Load("../../../.env")
}
