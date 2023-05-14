package main

import (
	"context"
	"log"
	"web-app-test/handlers/proxyhdl"
	"web-app-test/middleware"
	"web-app-test/repository/mongorepo"
	"web-app-test/repository/proxyrepo"
	"web-app-test/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func loadEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
}

func main() {
	loadEnv()
	mongoDbUrl, ok := viper.Get("MONGODB_URL").(string)
	if !ok {
		log.Fatalln("Invalid database URL")
	}
	mongoRepository := mongorepo.NewDB(mongoDbUrl)
	defer mongoRepository.Client.Disconnect(context.TODO())

	proxyHost, ok := viper.Get("PROXY_HOST").(string)
	if !ok {
		log.Fatalln("Invalid proxy host")
	}
	proxyRepository := proxyrepo.New(proxyHost)
	proxyService := service.New(proxyRepository, mongoRepository)
	proxyHandler := proxyhdl.NewHTTPHandler(proxyService)

	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	ratelimit, ok := viper.Get("RATELIMIT").(string)
	if !ok {
		log.Fatalln("Invalid ratelimit")
	}
	app.Use(middleware.RateLimiter(ratelimit))
	app.GET("*path", proxyHandler.Get)
	app.Run(":8080")
}
