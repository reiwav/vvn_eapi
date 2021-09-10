package ref_request

import (
	"eapi/dao"
	"eapi/reiway/tibero"

	"eapi/o/request"
)

var tableRefRequest = tibero.Table{
	TableName: "ref_requests",
	DB:        dao.Database(),
}

func NewTable() error {
	return tableRefRequest.NewTable(&RefRequest{})
}

func SelectMany(where map[string]string, order string, skip, limit int) ([]RefRequest, error) {
	var res = []RefRequest{}
	return res, tableRefRequest.SelectMany(where, order, skip, limit, &res)
}

func SelecOne(where map[string]string, order string) (*RefRequest, error) {
	var res = &RefRequest{}
	return res, tableRefRequest.SelectOne(where, order, 0, 1, &res)
}

func (r *RefRequest) Insert() error {
	return tableRefRequest.InsertOne(r)
}

type RefRequest struct {
	tibero.BaseModel `rei:"inline"`
	FromRequestID    int64           `rei:"form_request_id" json:"formRequestID"`
	RequestID        int64           `rei:"request_id" json:"requestID"`
	FromRequest      request.Request `rei:"-" json:"fromRequest"`
	Request          request.Request `rei:"-" json:"request"`
}
