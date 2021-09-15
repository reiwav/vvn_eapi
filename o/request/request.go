package request

import (
	"eapi/dao"
	"eapi/o/file"
	"eapi/reiway/tibero"
)

var tableRequest = tibero.Table{
	TableName: "requests",
	DB:        dao.Database(),
}

func NewTable() error {
	return tableRequest.NewTable(&Request{})
}

func SelectMany(where map[string]string, order string, skip, limit int) ([]Request, error) {
	var res = []Request{}
	return res, tableRequest.SelectMany(where, order, skip, limit, &res)
}

func SelectCustomMany(where string, order string, skip, limit int) ([]Request, error) {
	var res = []Request{}
	return res, tableRequest.SelectCustomMany(where, order, skip, limit, &res)
}

func SelectDistinct(where map[string]string, cols []string, res interface{}) error {
	return tableRequest.SelectDistinct(where, cols, res)
}

func SelecOne(where map[string]string, order string) (*Request, error) {
	var res = &Request{}
	var err = tableRequest.SelectOne(where, order, 0, 1, &res)
	if err != nil {
		return nil, err
	}
	return res, res.CheckID()
}

func (r *Request) Insert() error {
	return tableRequest.InsertOne(r)
}
func (r *Request) UpdateByID() error {
	return tableRequest.UpdateByID(r)
}

type Request struct {
	tibero.BaseModel     `rei:"inline"`
	RequestID            tibero.String `rei:"request_id,10000" json:"requestID"`
	RequestBody          tibero.String `rei:"request_body" json:"requestBody"`
	RequestType          tibero.String `rei:"request_type,10000" json:"requestType"`
	ResponseBody         tibero.String `rei:"response_body" json:"responseBody"`
	CardType             int           `rei:"card_type" json:"cardType"`
	Files                []file.File   `rei:"-" json:"images"`
	ContentType          tibero.String `rei:"content_type,10000" json:"contentType"`
	ProcessServer        tibero.String `rei:"process_server" json:"processServer"`
	ProcessTime          float64       `rei:"process_time" json:"processTime"`
	PrepareTime          float64       `rei:"prepare_time" json:"prepareTime"`
	StartedUnixTime      float64       `rei:"started_unix_time" json:"startedUnixTime"`
	ErrorMessage         tibero.String `rei:"error_message" json:"errorMessage"`
	StatusCode           int           `rei:"status_code" json:"statusCode"`
	SaveImage            int           `rei:"save_image" json:"saveImage"`
	ResponseBodyIdentify tibero.String `rei:"response_body_identify" json:"responseBodyIdentify"`
	ResponseBodyName     tibero.String `rei:"response_body_name,10000" json:"responseBodyName"`
	FakeCode             tibero.String `rei:"fake_code,5000" json:"fakeCode"`
}
