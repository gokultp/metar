package main

import (
	"os"

	"github.com/gokultp/metar/apis"
)

func main() {
	RedisURL := os.Getenv("REDIS_URL")
	Port := ":" + os.Getenv("PORT")

	api := apis.NewAPI(RedisURL)
	api.InitRoutes()
	api.Listen(Port)

}
