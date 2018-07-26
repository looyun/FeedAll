package controllers

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/looyun/feedall/models"
	"github.com/mmcdole/gofeed"
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
	username := c.Data["username"]
	user := models.User{}
	err := models.FindOne(models.Users, bson.M{"username": username}, &user)
	if err != nil {
		return nil, err
	}
	feeds := []models.Feed{}
	err = models.FindAll(models.Feeds, bson.M{"_id": bson.M{"$in": user.SubscribeFeedIDs}}, &feeds)
	if err != nil {
		return nil, err
	}
	return feeds, nil

}

func GetUserItems(c *macaron.Context) (interface{}, error) {
	username := c.Data["username"]
	user := models.User{}
	err := models.FindOne(models.Users, bson.M{"username": username}, &user)
	if err != nil {
		return nil, err
	}
	items := []models.Item{}
	err = models.FindAll(models.Items, bson.M{"feedID": bson.M{"$in": user.SubscribeFeedIDs}}, &items)
	if err != nil {
		return nil, err
	}
	return items, nil

}

func GetStarItems(c *macaron.Context) (interface{}, error) {
	username := c.Data["username"]
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

func Subscribe(c *macaron.Context) error {
	username := c.Data["username"]
	feedurl := c.Query("feedurl")
	fmt.Println(username)
	fmt.Println(feedurl)

	// parse feed url
	fb := gofeed.NewParser()
	origin_feed, err := fb.ParseURL(feedurl)
	if err != nil {
		fmt.Println("Parse err: ", err)
	}

	// if user subscribe this feed
	user := models.User{}
	err = models.FindOne(models.Users, bson.M{"username": username}, &user)
	if err != nil {
		return err
	}

	for _, fl := range user.SubscribeFeedLinks {
		if fl == origin_feed.Link {
			return nil
		}
	}
	bs_feed, err := bson.Marshal(origin_feed)
	if err != nil {
		fmt.Println(err)
	}
	feed := models.Feed{}
	err = bson.Unmarshal(bs_feed, &feed)
	if err != nil {
		fmt.Println(err)
	}
	if feed.FeedLink == "" {
		feed.FeedLink = feedurl
	}
	result := bson.M{}
	err = models.FindOne(models.Feeds,
		bson.M{"link": feed.Link},
		&result)
	var id bson.ObjectId
	if err != nil {
		if err == mgo.ErrNotFound {
			id = bson.NewObjectId()
			feed.ID = id
			err := models.Insert(models.Feeds, feed)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		id = result["_id"].(bson.ObjectId)
	}
	err = models.Update(models.Users,
		bson.M{"username": username},
		bson.M{"$addToSet": bson.M{"subscribeFeedLinks": feedurl}},
	)
	if err != nil {
		return err
	}
	err = models.Update(models.Users,
		bson.M{"username": username},
		bson.M{"$addToSet": bson.M{"subscribeFeedIDs": id}},
	)
	if err != nil {
		return err
	}
	return nil
}
