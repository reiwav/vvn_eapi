package init2

import (
	"eapi/common"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/reiwav/x/file"

	"github.com/gin-gonic/gin"
)

func init() {
	initConfig()
	initLog()
	initDB()
}

func initConfig() {
	// Open our jsonFile
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully Opened config.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var res = common.ConfigSystem{}
	json.Unmarshal([]byte(byteValue), &res)
	common.ConfigSystemCache = &res
}

func initDB() {
	var conf = common.ConfigSystemCache.DB
	_, err := conf.Connect(true)
	if err != nil {
		fmt.Println("ERROR CONNECT DB:", err)
		panic(err)
	}
}

func initLog() {
	file.CreateFolder("./log")
	file.CreateFolder(common.ConfigSystemCache.FolderAdmin)
	//config for gin request log
	{
		f, _ := os.Create(filepath.Join("log", "gin.log"))
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}
	//config for app log use glog
	{
		flag.Set("alsologtostderr", "true")
		flag.Set("log_dir", "./log")
		flag.Parse()
	}
}
