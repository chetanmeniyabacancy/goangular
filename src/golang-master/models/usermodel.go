package models

import (
	"github.com/jmoiron/sqlx"
	"golang-master/lang"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)
type User struct {
    Id	  	int64 `json:"id"`
    Email   string `json:"email"`
	Hash 	string `json:"hash"`
}

type ReqLogin struct {
    Email   string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Admin Login
func Login(db *sqlx.DB, reqlogin *ReqLogin)  (*User, string) {
	email := reqlogin.Email
	password := reqlogin.Password
	var user User

	err := db.Get(&user,"Select id,email,password as hash from admin_users where email = ? and password = ?",email,GetMD5Hash(password))
	fmt.Println(err)
	if err != nil {
		return &user, lang.Get("inavlid_username_or_password")
	}
	return &user, ""	
	
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}