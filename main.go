package main

import (
	// 1. init first
	_ "eapi/init2"

	//
	"eapi/common"
	"os"

	//2. middle
	"eapi/mid"
	//3. iniit 2nd
	"eapi/api"

	//"eapi/system"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	//"github.com/rs/cors"

	cors "github.com/rs/cors/wrapper/gin"
	//"net/http"
)

func main() {
	router := gin.New()

	router.Use(static.Serve("/", static.LocalFile("./admin", true)))
	router.Use(cors.AllowAll(), gin.Logger(), mid.Recovery())
	//api
	rootAPI := router.Group("/api")
	api.InitApi(rootAPI)
	//ws
	err := router.Run(common.ConfigSystemCache.Port)
	if err != nil {
		glog.Error(err)
		os.Exit(0)
	}
}
