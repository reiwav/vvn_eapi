package file

import (
	"eapi/dao"
	"eapi/reiway/tibero"
)

var tableFile = tibero.Table{
	TableName: "files",
	DB:        dao.Database(),
}

func NewTable() error {
	return tableFile.NewTable(&File{})
}

func SelectMany(where map[string]string, order string, skip, limit int) ([]File, error) {
	var res = []File{}
	return res, tableFile.SelectMany(where, order, skip, limit, &res)
}

func UnsafeSelectMany(cells, where, order string, skip, limit int) ([]File, error) {
	var res = []File{}
	return res, tableFile.UnsafeSelectMany(cells, where, order, skip, limit, &res)
}

func SelecOne(where map[string]string, order string) (*File, error) {
	var res = &File{}
	err := tableFile.SelectOne(where, order, 0, 0, &res)
	if err != nil {
		return nil, err
	}
	return res, res.CheckID()
}

func (r *File) Insert() error {
	return tableFile.InsertOne(r)
}

func InsertAll(fs []*File) error {
	var ts = tibero.LstModel{}
	for _, val := range fs {
		ts = append(ts, val)
	}
	return tableFile.InsertAll(ts)
}

func (r *File) UpdateByID() error {
	return tableFile.UpdateByID(r)
}

type File struct {
	tibero.BaseModel `rei:"inline"`
	RequestID        tibero.String `rei:"request_id,10000"  json:"requestID"`
	FilePath         tibero.String `rei:"file_path,10000"  json:"filePath"`
	FileRefer        int64         `rei:"file_refer"  json:"fileRefer"`
	Hash             tibero.String `rei:"hash,10000" json:"hash"`
	Type             tibero.String `rei:"type,5000" json:"type"`
	Version          tibero.String `rei:"version,5000"  json:"version"`
}
