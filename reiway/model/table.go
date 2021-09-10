package model

import (
	"eapi/o/file"
	"eapi/o/ref_request"
	"eapi/o/request"
	"fmt"

	"database/sql"
)

func SetTables(db *sql.DB) error {
	err := file.NewTable()
	if err != nil {
		return err
	}
	// // var req1 = file.File{
	// // 	RequestID: "request_2",
	// // 	FilePath:  "path 2",
	// // 	Version:   "2.2",
	// // }
	// // err = req1.Insert()
	// fmt.Println("======== CREATE VALUE: ", req1, err)
	var where = map[string]string{}
	//where["id"] = "1"
	res, err1 := file.SelectMany(where, "", 0, 0)
	fmt.Println("======== Select VALUE: ", res, err1)

	err = request.NewTable()
	if err != nil {
		return err
	}
	// var rq = request.Request{
	// 	RequestID:   "rq2",
	// 	RequestBody: "body2",
	// }
	// rq.Insert()
	where = map[string]string{}
	where["id"] = "1"
	res2, err2 := request.SelectMany(where, "", 0, 0)
	fmt.Println("======== Select VALUE: ", res2, err2)
	//err = ref_request.SetTable(db, true)

	err = ref_request.NewTable()
	if err != nil {
		return err
	}
	return nil
}
