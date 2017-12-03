package controllers

import (
	"github.com/looyun/feedall/models"

	"crypto/md5"
	"encoding/hex"
	"gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2/bson"
)

func Login(c *macaron.Context) bool {

	username := c.Query("username")
	password := c.Query("password")
	if username != "" && password != "" {
		user := models.User{}
		if models.GetUserInfo(models.Users, bson.M{"username": username}, &user) != true {
			c.Data["LoginErr"] = true
			return false
		} else {
			if user.Password == GetMd5String(password) {
				c.SetSecureCookie("username", username)
				c.SetSecureCookie("password", GetMd5String(password))
				c.Redirect("/")
				return true
			} else {
				c.Data["LoginErr"] = true
				return false
			}
		}
	}

	return false
}

func CheckLogin(c *macaron.Context) bool {
	username, u := c.GetSecureCookie("username")
	if u == false {
		return false
	}
	password, p := c.GetSecureCookie("password")
	if p == false {
		return false
	}
	user := models.User{}
	if username != "" {
		if models.GetUserInfo(models.Users, bson.M{"username": username}, &user) == true {
			return user.Password == password
		}
	}
	return false
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum([]byte("feedall")))
}
