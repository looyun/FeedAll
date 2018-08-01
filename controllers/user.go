package controllers

import (
	"errors"
	"fmt"
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
	page := c.QueryInt("page")
	if page > 0 {
		page--
	}
	perPage := c.QueryInt("per_page")
	if perPage == 0 {
		perPage = 30
	}
	if perPage > 100 {
		perPage = 100
	}
	username := c.Data["username"]
	user := models.User{}
	err := models.FindOne(models.Users, bson.M{"username": username}, &user)
	if err != nil {
		return nil, err
	}

	items := []models.Item{}
	err = models.Items.Find(bson.M{"feedID": bson.M{"$in": user.SubscribeFeedIDs}}).Sort("-publishedParsed").Skip(page * perPage).Limit(perPage).All(&items)
	if err != nil {
		return nil, err
	}
	return items, nil

}

func GetStarItems(c *macaron.Context) (interface{}, error) {
	page := c.QueryInt("page")
	if page > 0 {
		page--
	}
	perPage := c.QueryInt("per_page")
	if perPage == 0 {
		perPage = 30
	}
	if perPage > 100 {
		perPage = 100
	}
	username := c.Data["username"]
	user := models.User{}
	err := models.FindOne(models.Users, bson.M{"username": username}, &user)
	if err != nil {
		return nil, err
	}

	items := []models.Item{}
	err = models.Items.Find(bson.M{"_id": bson.M{"$in": user.StarItems}}).Sort("-publishedParsed").Skip(page * perPage).Limit(perPage).All(&items)
	if err != nil {
		return nil, err
	}
	return items, nil

}

func Subscribe(c *macaron.Context) error {
	username := c.Data["username"]
	feedurl := c.Query("feedurl")

	//tell if user subscribe this feed or not
	err := models.FindOne(models.Users,
		bson.M{
			"$and": []bson.M{
				bson.M{"username": username}, bson.M{"subscribeFeedURLs": feedurl}}},
		nil)

	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Println("not subscribe yet.")
		} else {
			return err
		}
	} else {
		fmt.Println("Already subscribe.")
		return nil
	}

	//tell if feed exist or not
	result := bson.M{}
	err = models.FindOne(
		models.Feeds,
		bson.M{"feedURLs": feedurl},
		&result)

	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Println("Not exist feed.")
		} else {
			return err
		}
	} else {
		fmt.Println("Exist feed.")
		err = models.Update(models.Users,
			bson.M{"username": username},
			bson.M{"$addToSet": bson.M{"subscribeFeedURLs": feedurl}},
		)
		if err != nil {
			return err
		}
		return nil
	}

	// Insert feed and update items
	err = InsertFeedAndUpdateItems(feedurl)
	if err != nil {
		return err
	}
	return nil
}
