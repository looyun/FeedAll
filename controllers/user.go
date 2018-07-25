package controllers

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/looyun/feedall/models"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var TokenSecure []byte = []byte("feedall")

func Signup(c *macaron.Context) error {
	username := c.Query("username")
	password := c.Query("password")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = models.FindOne(models.Users, bson.M{"username": username}, nil)

	if err == nil {
		return errors.New("user exist.")
	}
	if err != mgo.ErrNotFound {
		return err
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

func Login(c *macaron.Context) (string, error) {
	username := c.Query("username")
	password := c.Query("password")
	user := models.User{}
	err := models.FindOne(models.Users, bson.M{"username": username}, &user)

	if err != nil {
		return "", err
	}
	hash := user.Hash

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", err
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString(TokenSecure)
	if err != nil {
		return "", err
	}
	return t, nil
}

func GetUserFeeds(c *macaron.Context) (interface{}, error) {
	username := c.Params("username")
	user := models.User{}
	err := models.FindOne(models.Users, bson.M{"username": username}, &user)
	if err != nil {
		return nil, err
	}
	feeds := []models.Feed{}
	err = models.FindAll(models.Feeds, bson.M{"_id": bson.M{"$in": user.SubscribeFeedID}}, &feeds)
	if err != nil {
		return nil, err
	}
	return feeds, nil

}

func GetUserItems(c *macaron.Context) (interface{}, error) {
	username := c.Params("username")
	user := models.User{}
	err := models.FindOne(models.Users, bson.M{"username": username}, &user)
	if err != nil {
		return nil, err
	}
	items := []models.Item{}
	err = models.FindAll(models.Items, bson.M{"feedID": bson.M{"$in": user.SubscribeFeedID}}, &items)
	if err != nil {
		return nil, err
	}
	return items, nil

}

func GetStarItems(c *macaron.Context) (interface{}, error) {
	username := c.Params("username")
	user := models.User{}
	err := models.FindOne(models.Users, bson.M{"username": username}, &user)
	if err != nil {
		return nil, err
	}
	items := []models.Item{}
	err = models.FindAll(models.Items, bson.M{"_id": bson.M{"$in": user.StarItems}}, &items)
	if err != nil {
		return nil, err
	}
	return items, nil

}
