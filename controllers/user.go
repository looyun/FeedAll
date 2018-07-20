package controllers

import (
	"github.com/go-macaron/session"
	"github.com/looyun/feedall/models"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Signin(c *macaron.Context) error {
	username := c.Query("username")
	password := c.Query("password")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = models.FindOne(models.Users, bson.M{"username": username}, nil)

	if err != nil {
		if err != mgo.ErrNotFound {
			return err
		}
	}
	user := models.User{
		ID:       bson.NewObjectId(),
		Username: username,
		Hash:     string(hash),
	}
	err = models.Insert(models.Users, user)
	if err != nil {
		return err
	}
	return nil
}

func Login(c *macaron.Context, s session.Store) error {
	username := c.Query("username")
	password := c.Query("password")
	user := models.User{}
	err := models.FindOne(models.Users, bson.M{"username": username}, &user)

	if err != nil {
		return err
	}
	hash := user.Hash

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	sessionID := uuid.Must(uuid.NewV4())
	s.Set(sessionID, username)
	c.SetSecureCookie("user", username)
	c.SetSecureCookie("session", string(sessionID[:]))
	return nil
}

func CheckLogin(c *macaron.Context, s session.Store) bool {
	session, _ := c.GetSecureCookie("session")
	username, _ := c.GetSecureCookie("user")

	if s.Get(session) == username {
		return true
	}
	return false
}
