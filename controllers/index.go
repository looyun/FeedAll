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

var TelegramBotToken string = "123456"

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
			bson.M{"feedLink": bson.M{"$in": user.SubscribeFeedID}},
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
					bson.M{"$match": bson.M{"feedLink": bson.M{"$in": user.SubscribeFeedID}}},
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

// GetFeeds get feeds sort by subscribeCount.
func GetFeeds(c *macaron.Context) interface{} {

	page := c.QueryInt("page")
	if page == 0 {
		page = 1
	}
	perPage := c.QueryInt("perPage")
	if perPage == 0 {
		perPage = 20
	}
	feeds := []bson.M{}
	models.Feeds.Find(bson.M{}).Sort("-subscribeCount").Skip((page - 1) * perPage).Limit(perPage).All(&feeds)

	return feeds
}

// GetFeeds get feeds sort by subscribeCount.
func GetFeedItems(c *macaron.Context) interface{} {

	page := c.QueryInt("page")
	if page == 0 {
		page = 1
	}
	perPage := c.QueryInt("perPage")
	if perPage == 0 {
		perPage = 20
	}

	feedlink := ParseURL(c.Params(":feedlink"))
	feed := bson.M{}
	models.FindOne(models.Feeds,
		bson.M{"feedLink": feedlink},
		&feed)
	feedID := feed["_id"]

	items := []bson.M{}
	models.Items.Find(bson.M{"feedID": feedID}).Sort("-publishedParsed").Skip(perPage * (page - 1)).Limit(perPage).All(&items)

	return items
}
func GetFeed(c *macaron.Context) interface{} {
	feedlink := ParseURL(c.Params(":feedlink"))
	feed := bson.M{}
	models.FindOne(models.Feeds,
		bson.M{"feedLink": feedlink},
		&feed)

	return feed
}

func GetItems(c *macaron.Context, n int) interface{} {
	if n > 100 {
		return nil
	}
	items := []bson.M{}
	models.FindSortLimit(models.Items,
		bson.M{}, "-starCount", n,
		&items)
	return items
}
func GetItem(c *macaron.Context) interface{} {
	itemlink := ParseURL(c.Params(":itemlink"))
	item := bson.M{}
	models.FindOne(models.Items,
		bson.M{"link": itemlink},
		&item)
	return item

}

func GetRandomItem(c *macaron.Context, n int) interface{} {
	if n > 100 {
		return nil
	}
	items := []bson.M{}
	models.PipeAll(models.Items, []bson.M{{"$sample": bson.M{"size": n}}},
		&items)
	return items

}
func GetLatestItem(c *macaron.Context, n int) interface{} {
	if n > 100 {
		return nil
	}
	items := []bson.M{}
	models.FindSortLimit(models.Items,
		bson.M{}, "-publishedParsed", n,
		&items)
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
	if err != nil {
		fmt.Println(err)
	}
	return u
}
