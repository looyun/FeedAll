package controllers

import (
	"feedall/models"
	"fmt"
	"gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2/bson"
)

func Register(c *macaron.Context) bool {
	username := c.Query("username")
	password := c.Query("password")
	if username != "" && password != "" {
		user := models.User{
			ID:       bson.NewObjectId(),
			Username: username,
			Password: GetMd5String(password),
		}
		if !models.GetUserInfo(models.Users, bson.M{"username": username}, &user) {
			if models.Insert(models.Users, user) {
				fmt.Println("register successful!")
				c.Redirect("/user/login")
				return true
			} else {
				fmt.Println("register failed!")
				c.Redirect("/user/register")
			}
		} else {
			fmt.Println("already register!")
			c.Redirect("/user/login")
		}
	}
	return false
}
