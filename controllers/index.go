package controllers

import (
	// "encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/looyun/feedall/models"

	"gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2/bson"
)

func GetUserFeed(c *macaron.Context) {
	if !CheckLogin(c) {
		c.HTML(200, "login")

		return
	}
	username, _ := c.GetSecureCookie("username")
	user := models.User{}
	feed := make([]*models.Feed, 0)
	item := []bson.M{}
	if models.FindOne(models.Users, bson.M{"username": username}, &user) == true {

		fmt.Println("parse ", user.Username, " feed!")
		models.FindSort(models.Feeds,
			bson.M{"feedLink": bson.M{"$in": user.FeedLink}},
			"items",
			&feed)
		if len(feed) == 0 {
			c.Data["Hello"] = true
			fmt.Println("feeds ", "no match")
		} else {
			fmt.Println("feeds ", "match")
		}
		if feedlink := c.Query("feedlink"); feedlink != "" {

			fmt.Println("hi", c.Query("feedlink"))
			feedlink = ParseURL(feedlink)

			models.PipeAll(models.Feeds,
				[]bson.M{
					bson.M{"$match": bson.M{"feedLink": feedlink}},
					bson.M{"$unwind": "$items"},
					bson.M{"$sort": bson.M{"items.publishedParsed": -1}},
					bson.M{"$limit": 45},
				},
				&item)
			if len(item) == 0 {
				fmt.Println("items ", "no match")
			} else {
				fmt.Println("items ", "match")
			}
		} else {
			c.Data["root"] = true
			models.PipeAll(models.Feeds,
				[]bson.M{
					bson.M{"$match": bson.M{"feedLink": bson.M{"$in": user.FeedLink}}},
					bson.M{"$unwind": "$items"},
					bson.M{"$sort": bson.M{"items.publishedParsed": -1}},
					bson.M{"$limit": 45},
				},
				&item)
			if len(item) == 0 {
				fmt.Println("items ", "no match")
			} else {
				fmt.Println("items ", "match")
			}
		}
		c.Data["User"] = user
		c.Data["Feed"] = feed
		c.Data["Item"] = item
	}
}

func GetFeed(c *macaron.Context) interface{} {
	feedlink := ParseURL(c.Params("*"))
	feed := bson.M{}
	models.FindOne(models.Feeds,
		bson.M{"feedLink": feedlink},
		&feed)

	return feed
}

func GetItem(c *macaron.Context) interface{} {
	itemlink := ParseURL(c.Params("*"))
	item := bson.M{}
	models.PipeOne(models.Items,
		bson.M{"link": itemlink},
		&item)
	fmt.Println(item)
	return item

}

func GetRandomItem(c *macaron.Context, n int) interface{} {
	items := []bson.M{}
	models.PipeAll(models.Items, []bson.M{{"$sample": bson.M{"size": n}}},
		&items)
	fmt.Println(items)
	return items

}
func GetLatestItem(c *macaron.Context, n int) interface{} {
	items := []bson.M{}
	models.FindSortLimit(models.Items,
		bson.M{}, "-publishedParsed", n,
		&items)
	fmt.Println(items)
	return items

}

func StandarURL(s string) string {
	if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}
	return s
}

func ParseURL(s string) string {
	u, err := url.QueryUnescape(s)
	fmt.Println(err)
	return u
}
