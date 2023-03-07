package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"paopao-ce-teaching/internal/conf"
	"paopao-ce-teaching/internal/routers"
)

func init() {
	conf.Initialize()
}

func main() {
	gin.SetMode(conf.ServerSetting.RunMode)
	router := routers.NewRouter()

	if err := router.Run(conf.ServerSetting.HttpIp + ":" + conf.ServerSetting.HttpPort); err != nil {
		log.Fatalf("run app failed: %s", err)
	}
}
