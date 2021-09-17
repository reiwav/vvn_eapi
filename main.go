package main

import (
	// 1. init first
	_ "eapi/init2"
	"eapi/mid"

	//
	"eapi/common"
	"os"

	//2. middle

	//3. iniit 2nd
	"eapi/api"

	//"eapi/system"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	cors "github.com/rs/cors/wrapper/gin"
	//"github.com/rs/cors"
	//"net/http"
)

func main() {
	router := gin.New()

	router.Use(common.StaticServe("./admin"))
	router.Use(cors.AllowAll(), gin.Logger(), mid.Recovery())
	//api
	rootAPI := router.Group("/api")
	api.InitApi(rootAPI, router)

	//ws
	err := router.Run(common.ConfigSystemCache.Port)
	if err != nil {
		glog.Error(err)
		os.Exit(0)
	}
}
