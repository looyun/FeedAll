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
