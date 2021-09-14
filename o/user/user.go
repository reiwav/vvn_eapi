package user

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"eapi/dao"
	"eapi/reiway/tibero"
	"encoding/hex"
	"fmt"
	"io"
)

var tableUser = &tibero.Table{
	TableName: "jhi_user",
	DB:        dao.Database(),
}

func NewTable() error {
	return tableUser.NewTable(&User{})
}

func GetTable() *tibero.Table {
	return tableUser
}

// type User struct {
// 	tibero.BaseModel `rei:",inline"`
// 	Fullname         string `rei:"fullname" json:"fullname"`
// 	Birthday         string `rei:"birthday" json:"birthday"`
// 	Email            string `rei:"email" json:"email"`
// 	Type             string `rei:"type" json:"type"`
// 	Username         string `rei:"username" json:"username"`
// 	Password         string `rei:"password" json:"password"`
// }

type User struct {
	tibero.BaseModel `rei:"inline"`
	Login            tibero.String `json:"login" rei:"login"`
	FirstName        tibero.String `json:"first_name" rei:"firstName"`
	LastName         tibero.String `json:"last_name" rei:"lastName"`
	Email            tibero.String `rei:"email" json:"email" `
	ImageURL         tibero.String `rei:"image_url" json:"imageUrl"`
	Activated        tibero.Bool   `rei:"activated" json:"activated"`
	LangKey          tibero.String `rei:"lang_key" json:"langKey"`
	CreatedBy        tibero.String `rei:"created_by" json:"createdBy"`
	//CreatedDate      interface{}   `rei:"email" json:"createdDate"`
	LastModifiedBy tibero.String `rei:"lastModified_by" json:"lastModifiedBy"`
	//LastModifiedDate time.Time     `rei:"email" json:"lastModifiedDate"`
	//Authorities      []interface{} `rei:"email" json:"authorities"`
	MinioEndpoint tibero.String `rei:"minio_endpoint" json:"minioEndpoint"`
	MinioKey      tibero.String `rei:"minio_key" json:"minioKey"`
	MinioSecret   tibero.String `rei:"minio_secret" json:"minioSecret"`
	MinioUseSSL   tibero.Bool   `rei:"minio_use_ssl" json:"minioUseSSL"`
	MinioBucket   tibero.String `rei:"minio_bucket" json:"minioBucket"`
	MinioPrefix   tibero.String `rei:"minio_prefix" json:"minioPrefix"`
	Password      tibero.String `rei:"password" json:"password"`
}

const Key = "f016f1185acaf153fab0bc4449803b5c"

func Encrypt(stringToEncrypt string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(Key)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func GetByUserPass(usr, pwd string) (*User, error) {
	var u = User{}
	var where = make(map[string]string)
	where["login"] = "'" + usr + "'"
	//where["password"] = "'" + Encrypt(pwd) + "'"
	return &u, tableUser.SelectOne(where, "", 0, 0, &u)
}

func SelectUser(usr, pwd string) (*User, error) {
	var where = make(map[string]string)
	where["login"] = "'" + usr + "'"
	//where["password"] = "'" + Encrypt(pwd) + "'"
	var u = User{}
	var rows, err = tableUser.SelectRows(where, nil, "", 0, 0)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return &u, u.mapValue(rows)
}

func (u *User) mapValue(rows *sql.Rows) error {
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.Login, &u.FirstName,
			&u.LastName, &u.Email, &u.ImageURL, &u.Activated, &u.LangKey, &u.CreatedBy, &u.LastModifiedBy,
			&u.MinioEndpoint, &u.MinioKey, &u.MinioSecret, &u.MinioUseSSL, &u.MinioBucket, &u.MinioPrefix, &u.Password)
		if err != nil {
			return err
		}
		rows.Close()
	}
	return nil
}

func mapValues(rows *sql.Rows) ([]*User, error) {
	var us = make([]*User, 0)
	for rows.Next() {
		var u = User{}
		err := rows.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.Login, &u.FirstName,
			&u.LastName, &u.Email, &u.ImageURL, &u.Activated, &u.LangKey, &u.CreatedBy, &u.LastModifiedBy,
			&u.MinioEndpoint, &u.MinioKey, &u.MinioSecret, &u.MinioUseSSL, &u.MinioBucket, &u.MinioPrefix, &u.Password)
		if err != nil {
			return nil, err
		}
		us = append(us, &u)
	}
	return us, nil
}

func (u *User) Create() error {
	//pwd, _ := auth.GererateHashedPassword(string(u.Password))
	pwd := Encrypt(string(u.Password))
	u.Password = tibero.String(pwd)
	err := tableUser.InsertOne(u)
	return err
}

func (u *User) Update(isUpPass bool) error {
	if u.Password != "" && isUpPass {
		//pwd, _ := auth.GererateHashedPassword(string(u.Password))
		pwd := Encrypt(string(u.Password))
		u.Password = tibero.String(pwd)
	}
	err := tableUser.UpdateByID(u)
	return err
}

func GetAll() ([]User, error) {
	var urs = []User{}
	err := tableUser.SelectMany(nil, "", 0, 0, &urs)
	return urs, err
}

func (u *User) MarkDelete() error {
	err := tableUser.DeleteByID(u)
	return err
}

func GetByID(uID string) (*User, error) {
	var u = User{}
	where := make(map[string]string)
	where["id"] = uID
	var err = tableUser.SelectOne(where, "", 0, 0, &u)
	return &u, err
}
