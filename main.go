package main

import (
	"os"

	"github.com/gokultp/metar/apis"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	redisPassword := "" //default
	redisDB := 0        //default db

	port := ":" + os.Getenv("PORT")

	api := apis.NewAPI(redisURL, redisPassword, redisDB)
	api.InitRoutes()
	api.Listen(port)

}
