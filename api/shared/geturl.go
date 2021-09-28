package shared

import (
	"eapi/common"
	"eapi/dao"
	"eapi/o/file"
	"eapi/o/user"
	"eapi/reiway/tibero"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go"

	"github.com/gin-gonic/gin"
)

func getUrlImage(endpoint, bucket, pathFile string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}
	u.Path = path.Join(u.Path, bucket, pathFile)
	s := u.String()
	return s
}

func GetUrlFilePath(ctx *gin.Context, usr *user.User, pathFile string) string {
	// var storage, err1 = getMiniStorage(usr.MinioEndpoint.String(), string(usr.MinioKey), string(usr.MinioSecret), string(usr.MinioBucket), bool(usr.MinioUseSSL))
	// rest.AssertNil(err1)
	// url := "N/A"
	// if files.FilePath != "" {
	// 	url, err = storage.GetURL(ctx, string(files.FilePath), time.Duration(common.ConfigSystemCache.TimeExpires)*time.Second)
	// 	fmt.Println("==============", url)
	// 	if err != nil {
	// 		rest.AssertNil(fmt.Errorf("Gen url %v \n", err))
	// 	}
	// 	files.FilePath = tibero.String(url)
	// }
	var cacheConf = common.ConfigSystemCache
	if cacheConf.UseMinio {
		var storage, _ = getMiniStorage(usr.MinioEndpoint.String(), string(usr.MinioKey), string(usr.MinioSecret), string(usr.MinioBucket), bool(usr.MinioUseSSL))
		if pathFile != "" && storage != nil {
			url, _ := storage.GetURL(ctx, string(pathFile), time.Duration(cacheConf.TimeExpires)*time.Second)
			fmt.Println("==============", url)
			return url
		}
	}
	return getUrlImage(string(usr.MinioEndpoint), string(usr.MinioBucket), pathFile)
}

func getMiniStorage(endpoint, minioKey, minioSecret, minioBucket string, minioUseSSL bool) (*dao.MinioStorage, error) {
	if strings.Contains(string(endpoint), "https://") {
		endpoint = strings.ReplaceAll(endpoint, "https://", "")
	} else if strings.Contains(string(endpoint), "http://") {
		endpoint = strings.ReplaceAll(endpoint, "http://", "")
	}
	minioClient, err := minio.New(endpoint, minioKey, minioSecret, minioUseSSL)
	if err != nil {
		return nil, err
	}
	var storage = dao.NewMinioStorage(minioClient, minioBucket)
	return storage, nil
}

func GetUrlImageArrFiles(ctx *gin.Context, files []file.File, usr *user.User) []file.File {
	var fileNews = make([]file.File, len(files))
	var cacheConf = common.ConfigSystemCache
	var storage *dao.MinioStorage
	if cacheConf.UseMinio {
		storage, _ = getMiniStorage(usr.MinioEndpoint.String(), string(usr.MinioKey), string(usr.MinioSecret), string(usr.MinioBucket), bool(usr.MinioUseSSL))
	}
	for i, val := range files {
		if cacheConf.UseMinio {
			if val.FilePath != "" && storage != nil {
				url, _ := storage.GetURL(ctx, string(val.FilePath), time.Duration(cacheConf.TimeExpires)*time.Second)
				val.FilePath = tibero.String(url)
			}
		} else {
			val.FilePath = tibero.String(getUrlImage(string(usr.MinioEndpoint), string(usr.MinioBucket), val.FilePath.String()))
		}
		fileNews[i] = val
	}
	return fileNews
}
