package api

import (
	"eapi/api/admin"
	"eapi/dao/settable"
	"eapi/mid"
	"eapi/o/file"
	"eapi/o/request"
	"eapi/o/user"
	"eapi/reiway/tibero"
	"fmt"
	"strings"
	"time"

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

func InitApi(root *gin.RouterGroup) {
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
	var rqID = ctx.Query("request_id")
	var where = make(map[string]string)
	where["id"] = rqID
	var req, err = request.SelecOne(where, "")
	rest.AssertNil(err)
	where = make(map[string]string)
	where["request_id"] = rqID
	var files, _ = file.SelectMany(where, " id desc", 0, 0)
	req.Files = files
	js.SendString(ctx, req)
}

func handleGetImage(ctx *gin.Context) {
	var rqID = ctx.Query("image_id")
	var where = make(map[string]string)
	where["id"] = rqID
	var files, _ = file.SelecOne(where, "")
	js.SendString(ctx, files)
}

func handleAccount(ctx *gin.Context) {
	var u = mid.GetMyUser(ctx)
	js.SendString(ctx, u)
}

func handlePostAccount(ctx *gin.Context) {
	var u *user.User
	rest.AssertNil(ctx.BindJSON(&u))
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
	if len(sorts) > 0 {
		for _, val := range sorts {
			var sort1s = strings.Split(val, ",")
			if len(sorts) > 0 {
				oderBy = sort1s[0] + " " + sort1s[1]
			} else {
				oderBy = sort1s[0] + " asc"
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
		where += " AND created_at between TO_DATE('" + start + "') AND TO_DATE('" + end + "')"
	}
	var res, err = request.SelectCustomMany(where, oderBy, int(page*size), int(size))
	fmt.Println("======== REQUESTS", err)
	js.SendString(ctx, res)
}

func handleImages(ctx *gin.Context) {
	var q = ctx.Request.URL.Query()
	var page = web.MustGetInt64("page", q)
	var size = web.MustGetInt64("size", q)
	var sorts = strings.Split(ctx.Query("sort"), ",")
	var oderBy = sorts[0] + " " + sorts[1]
	var res, _ = file.SelectMany(nil, oderBy, int(page*size), int(size))
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
