package controllers

import (
	// "encoding/json"
	"feedall/models"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

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
	fmt.Println("hi", c.Query("feedlink"))
	fmt.Println("hi", c.Query("page"))
	if models.GetUserInfo(models.Users, bson.M{"username": username}, &user) == true {

		if c.Query("page") != "" {
			GetMoreItem(c, user.FeedURL)
			return
		}

		fmt.Println("parse ", user.Username, " feed!")
		models.GetFeed(models.Feeds,
			bson.M{"feedLink": bson.M{"$in": user.FeedURL}},
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

			models.GetAllItem(models.Feeds,
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
			models.GetAllItem(models.Feeds,
				[]bson.M{
					bson.M{"$match": bson.M{"feedLink": bson.M{"$in": user.FeedURL}}},
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

		c.HTML(200, "index")
	}
}
func GetMoreItem(c *macaron.Context, s []string) {

	item := []bson.M{}
	page, _ := strconv.Atoi(c.Query("page"))
	feedlink := c.Query("feedlink")
	if feedlink != "" {
		feedlink = ParseURL(feedlink)
		models.GetAllItem(models.Feeds,
			[]bson.M{
				bson.M{"$match": bson.M{"feedLink": feedlink}},
				bson.M{"$unwind": "$items"},
				bson.M{"$sort": bson.M{"items.publishedParsed": -1}},
				bson.M{"$skip": 45 + page*30},
				bson.M{"$limit": 30},
			},
			&item)
	} else {
		models.GetAllItem(models.Feeds,
			[]bson.M{
				bson.M{"$match": bson.M{"feedLink": bson.M{"$in": s}}},
				bson.M{"$unwind": "$items"},
				bson.M{"$sort": bson.M{"items.publishedParsed": -1}},
				bson.M{"$skip": 45 + page*30},
				bson.M{"$limit": 30},
			},
			&item)
	}
	c.Data["ITEMS"] = item
	c.HTML(200, "items")
	return
}

func GetItemContent(c *macaron.Context) {
	itemlink := ParseURL(c.Params("*"))
	fmt.Println(itemlink)
	feed := bson.M{}
	models.GetItem(models.Feeds,
		[]bson.M{
			bson.M{"$match": bson.M{"items.link": itemlink}},
			bson.M{"$unwind": "$items"},
			bson.M{"$match": bson.M{"items.link": itemlink}},
		},
		&feed)

	c.Data["Feed"] = feed
	item, _ := feed["items"].(bson.M)
	content, _ := item["content"].(string)
	io.WriteString(c, content)

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
