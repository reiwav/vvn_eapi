package api

import (
	"eapi/api/admin"
	"eapi/api/shared"
	"eapi/common"
	"eapi/dao"
	"eapi/dao/settable"
	"eapi/mid"
	"eapi/o/file"
	"eapi/o/request"
	"eapi/o/user"
	"eapi/reiway/tibero"
	"fmt"
	"strings"
	"time"

	minio "github.com/minio/minio-go"

	"eapi/api/auth"
	apiRq "eapi/api/request"
	eAuth "eapi/ctrl/auth"

	"github.com/reiwav/x/mlog"
	"github.com/reiwav/x/rest"
	"github.com/reiwav/x/web"

	"github.com/gin-gonic/gin"
)

var logData = mlog.NewTagLog("Done-Payment")

var js = rest.JsonRender{}

func InitApi(root *gin.RouterGroup, router *gin.Engine) {
	settable.CreateAllTable()
	admin.NewAdminServer(root, "admin")
	auth.NewAuthenServer(root, "account")
	apiRq.NewRequestServer(root, "common")
	root.POST("/authenticate", handleAuthenticate)
	root.POST("/register", handleRegister)
	root.Use(mid.AuthBasicJwt("", true))
	root.GET("/account", handleAccount)
	root.POST("/account", handlePostAccount)
	root.GET("/requests", handleRequests)
	root.GET("/requests/:request_id", handleGetRequest)
	root.GET("/images/:image_id", handleGetImage)
	//router.Use(mid.AuthBasicJwt("", true)).GET("/image/:id/view", HandleGetImageView)
	root.GET("/images", handleImages)

}

func handleRegister(ctx *gin.Context) {
	var body *user.User
	rest.AssertNil(ctx.BindJSON(&body))
	rest.AssertNil(body.Create())
	js.SendString(ctx, nil)
}

func handleAuthenticate(ctx *gin.Context) {
	var body = eAuth.LoginApp{}
	rest.AssertNil(ctx.BindJSON(&body))
	var tk, err = body.LoginUser()
	rest.AssertNil(err)
	js.SendString(ctx, map[string]string{
		"id_token": tk.MyToken.Token,
	})
}

func handleGetRequest(ctx *gin.Context) {
	var rqID = ctx.Params.ByName("request_id")
	fmt.Println("ID: ", rqID)
	var where = make(map[string]string)
	where["id"] = rqID
	var req, err = request.SelecOne(where, "")
	rest.AssertNil(err)
	var usr = mid.GetMyUser(ctx)
	where = make(map[string]string)
	where["file_refer"] = rqID
	var files, _ = file.SelectMany(where, " id desc", 0, 0)
	var fileNews = make([]file.File, len(files))
	var storage, err1 = getMiniStorage(usr.MinioEndpoint.String(), string(usr.MinioKey), string(usr.MinioSecret), string(usr.MinioBucket), bool(usr.MinioUseSSL))
	rest.AssertNil(err1)
	for i, val := range files {
		url := "N/A"
		if val.FilePath != "" {
			url, err = storage.GetURL(ctx, string(val.FilePath), time.Duration(common.ConfigSystemCache.TimeExpires)*time.Second)
			fmt.Println("==============", url)
			if err != nil {
				rest.AssertNil(fmt.Errorf("Gen url %v \n", err))
			}
			val.FilePath = tibero.String(url)
		}

		fileNews[i] = val
	}
	req.Files = fileNews
	js.SendString(ctx, req)
}

func handleGetImage(ctx *gin.Context) {
	var rqID = ctx.Params.ByName("image_id")
	var where = make(map[string]string)
	where["id"] = rqID
	var files, err = file.SelecOne(where, "")
	rest.AssertNil(err)
	if files != nil {
		var usr = mid.GetMyUser(ctx)
		var storage, err1 = getMiniStorage(usr.MinioEndpoint.String(), string(usr.MinioKey), string(usr.MinioSecret), string(usr.MinioBucket), bool(usr.MinioUseSSL))
		rest.AssertNil(err1)
		url := "N/A"
		if files.FilePath != "" {
			url, err = storage.GetURL(ctx, string(files.FilePath), time.Duration(common.ConfigSystemCache.TimeExpires)*time.Second)
			fmt.Println("==============", url)
			if err != nil {
				rest.AssertNil(fmt.Errorf("Gen url %v \n", err))
			}
			files.FilePath = tibero.String(url)
		}
	}
	js.SendString(ctx, files)
}

func HandleGetImageView(ctx *gin.Context) {
	fmt.Println(ctx.Params)
	var rqID = ctx.Params.ByName("id")
	var where = make(map[string]string)
	where["id"] = rqID
	var files, err = file.SelecOne(where, "")

	rest.AssertNil(err)

	if files != nil {
		var usr = mid.GetMyUser(ctx)
		var storage, err1 = getMiniStorage(usr.MinioEndpoint.String(), string(usr.MinioKey), string(usr.MinioSecret), string(usr.MinioBucket), bool(usr.MinioUseSSL))
		rest.AssertNil(err1)
		url := "N/A"
		if files.FilePath != "" {
			url, err = storage.GetURL(ctx, string(files.FilePath), time.Duration(common.ConfigSystemCache.TimeExpires)*time.Second)
			fmt.Println("==============", url)
			if err != nil {
				rest.AssertNil(fmt.Errorf("Gen url %v \n", err))
			}
			files.FilePath = tibero.String(url)
		}
	}
	js.SendString(ctx, files)
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

func handleAccount(ctx *gin.Context) {
	var u = mid.GetMyUser(ctx)
	js.SendString(ctx, u)
}

func handlePostAccount(ctx *gin.Context) {
	var u *user.User
	rest.AssertNil(ctx.BindJSON(&u))
	u.Password = ""
	err := u.Update(false)
	if err != nil {
		rest.AssertNil(rest.BadRequest(err.Error()))
	}
	js.SendString(ctx, nil)
}

func handleRequests(ctx *gin.Context) {
	var requestID = ctx.Query("requestID.contains")
	var requestType = ctx.Query("requestType.equals")
	var fakeCode = ctx.Query("fakeCode.equals")
	var createdAtGt = ctx.Query("createdAt.greaterThan")
	var createdAtLt = ctx.Query("createdAt.lessThan")
	var responseBodyIdentify = ctx.Query("responseBodyIdentify.contains")
	var responseBodyName = ctx.Query("responseBodyName.contains")
	var q = ctx.Request.URL.Query()
	var page = web.MustGetInt64("page", q)
	var size = web.MustGetInt64("size", q)
	var sorts = ctx.QueryArray("sort")
	var oderBy string
	var req = request.Request{}
	if len(sorts) > 0 {
		for _, val := range sorts {
			var ordField string
			if strings.Contains(val, ",") {
				var sort1s = strings.Split(val, ",")
				fmt.Println(sort1s)
				var fieldJson = sort1s[0]
				colTibero := shared.GetFieldNameTibero(fieldJson, req)
				ordField = colTibero + " " + sort1s[1]
			} else {
				ordField = val + " asc"
			}
			if oderBy != "" {
				oderBy += ", " + ordField
			} else {
				oderBy += ordField
			}
		}

	}

	var where = " deleted_at is null "
	if requestType != "" {
		where += " AND request_type = '" + requestType + "'"
	}
	if fakeCode != "" {
		where += " AND fake_code = '" + fakeCode + "'"
	}
	if requestID != "" {
		where += " AND request_id LIKE '%" + requestID + "%'"
	}
	if responseBodyName != "" {
		where += " AND response_body_name LIKE '%" + requestID + "%'"
	}
	if responseBodyIdentify != "" {
		where += " AND response_body_identify LIKE '%" + responseBodyIdentify + "%'"
	}

	if createdAtGt != "" && createdAtLt != "" {
		var _, start = ParseTimeQuery(createdAtGt)
		var _, end = ParseTimeQuery(createdAtLt)
		where += " AND created_at between " + start + " AND " + end
	}
	if page < 0 {
		page = 0
	}
	var skip = int(page * size)
	var res, err = request.SelectCustomSkipLimit(where, oderBy, skip, int((page+1)*size))
	rest.AssertNil(err)
	var total, _ = request.Count(nil)
	ctx.Writer.Header().Set("X-Total-Count", fmt.Sprintf("%v", total))
	js.SendString(ctx, res)
}

func handleImages(ctx *gin.Context) {
	var q = ctx.Request.URL.Query()
	var page = web.MustGetInt64("page", q)
	var size = web.MustGetInt64("size", q)
	var sorts = strings.Split(ctx.Query("sort"), ",")

	if page < 0 {
		page = 0
	}
	var skip = int(page * size)
	var oderBy = sorts[0] + " " + sorts[1]
	var res, _ = file.SelectSkipLimit(nil, oderBy, skip, int((page+1)*size))
	if len(res) > 0 {
		var usr = mid.GetMyUser(ctx)
		var storage, err1 = getMiniStorage(usr.MinioEndpoint.String(), string(usr.MinioKey), string(usr.MinioSecret), string(usr.MinioBucket), bool(usr.MinioUseSSL))
		rest.AssertNil(err1)

		var fileNews = make([]file.File, len(res))
		for i, val := range res {
			if val.FilePath != "" {
				url, err := storage.GetURL(ctx, string(val.FilePath), time.Duration(common.ConfigSystemCache.TimeExpires)*time.Second)
				if err != nil {
					rest.AssertNil(fmt.Errorf("Gen url %v \n", err))
				}
				val.FilePath = tibero.String(url)
			}
			fileNews[i] = val
		}
		res = fileNews
	}
	var total, _ = request.Count(nil)
	ctx.Writer.Header().Set("X-Total-Count", fmt.Sprintf("%v", total))
	js.SendString(ctx, res)
}

func ParseTimeQuery(d string) (time.Time, string) {
	t, err := time.Parse(time.RFC3339, d)
	if err != nil {
		layout1 := "2006-01-02T15:04:05.000Z"
		t, _ = time.Parse(layout1, d)
	}
	return t, tibero.ConvertTimeToDate(t)
}
