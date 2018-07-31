package controllers

import (

	// "encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/looyun/feedall/models"
	"github.com/mmcdole/gofeed"

	"gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2/bson"
)

var TelegramBotToken string = "123456"

func ParseURL(feedurl string) (interface{}, error) {

	// parse feed url
	fp := gofeed.NewParser()
	parsed_feed, err := fp.ParseURL(feedurl)
	if err != nil {
		return nil, err
	}

	bs_feed, err := bson.Marshal(parsed_feed)
	if err != nil {
		return nil, err
	}
	feed := models.Feed{}
	err = bson.Unmarshal(bs_feed, &feed)
	if err != nil {
		return nil, err
	}

	feed.FeedURL = feedurl
	feed.ID = bson.NewObjectId()
	return feed, nil
}

func UpdateItems(feed models.Feed) error {
	feedurl := feed.FeedURL

	// parse feed url
	fp := gofeed.NewParser()
	parsed_feed, err := fp.ParseURL(feedurl)
	if err != nil {
		return err
	}

	bs_items, err := bson.Marshal(parsed_feed.Items)
	if err != nil {
		return err
	}

	items := []models.Item{}
	err = bson.Unmarshal(bs_items, &items)
	if err != nil {
		return err
	}
	for _, item := range items {
		item.FeedID = feed.ID
		_, err := models.Upsert(models.Items,
			bson.M{
				"$and": []bson.M{
					bson.M{"feedID": feed.ID}, bson.M{"link": item.Link}}},
			item)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFeeds get feeds sort by subscribeCount.
func GetFeeds(c *macaron.Context) interface{} {

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
	feeds := []bson.M{}
	err := models.Feeds.Find(bson.M{}).Sort("-subscribeCount").Skip(page * perPage).Limit(perPage).All(&feeds)

	if err != nil {
		fmt.Println(err)
	}
	return feeds
}

// GetFeeds get feeds sort by subscribeCount.
func GetFeedItems(c *macaron.Context) interface{} {

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

	feedID := c.Params(":id")

	items := []bson.M{}
	err := models.Items.Find(bson.M{"feedID": feedID}).Sort("-publishedParsed").Skip(perPage * page).Limit(perPage).All(&items)

	if err != nil {
		fmt.Println(err)
	}
	return items
}
func GetFeed(c *macaron.Context) interface{} {
	id := c.Params(":id")
	feed := bson.M{}
	err := models.FindOne(models.Feeds,
		bson.M{"_id": id},
		&feed)
	if err != nil {
		fmt.Println(err)
	}

	return feed
}

func GetItems(c *macaron.Context, n int) interface{} {
	if n > 100 {
		return nil
	}
	items := []bson.M{}
	err := models.FindSortLimit(models.Items,
		bson.M{}, "-starCount", n,
		&items)
	if err != nil {
		fmt.Println(err)
	}
	return items
}
func GetItem(c *macaron.Context) interface{} {
	id := c.Params(":id")
	item := bson.M{}
	err := models.FindOne(models.Items,
		bson.M{"_id": id},
		&item)
	if err != nil {
		fmt.Println(err)
	}
	return item

}

func GetRandomItem(c *macaron.Context, n int) interface{} {
	if n > 100 {
		return nil
	}
	items := []bson.M{}
	err := models.PipeAll(models.Items, []bson.M{{"$sample": bson.M{"size": n}}},
		&items)
	if err != nil {
		fmt.Println(err)
	}
	return items

}
func GetLatestItem(c *macaron.Context, n int) interface{} {
	if n > 100 {
		return nil
	}
	items := []bson.M{}
	err := models.FindSortLimit(models.Items,
		bson.M{}, "-publishedParsed", n,
		&items)
	if err != nil {
		fmt.Println(err)
	}
	return items

}

func StandarURL(s string) string {
	if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}
	return s
}

func DecodeImg(str string, link string) string {
	str = strings.Replace(str, "&#34;", "\"", -1)
	str = strings.Replace(str, "src=\"/", "src=\""+link+"/", -1)
	return str
}

func ParseDate(t string) (then time.Time) {

	if len(t) >= 25 {
		if strings.HasSuffix(t, "0000") {
			then, _ = time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", t)
		} else if strings.HasSuffix(t, "GMT") {
			then, _ = time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", t)
		} else if strings.HasSuffix(t, "UTC") {
			then, _ = time.Parse("Mon, 02 Jan 2006 15:04:05 UTC", t)
		} else if strings.HasSuffix(t, "CST") {
			then, _ = time.Parse("Mon, 02 Jan 2006 15:04:05 CST", t)
		} else if strings.HasSuffix(t, "0400") {
			then, _ = time.Parse("Mon, 02 Jan 2006 15:04:05 -0400", t)
		} else if strings.HasSuffix(t, "Z") {
			then, _ = time.Parse(time.RFC3339, t)
		} else if strings.HasSuffix(t, "0800") {
			then, _ = time.Parse("Mon, 02 Jan 2006 15:04:05 +0800", t)
		}
	} else {
		if strings.HasSuffix(t, "0000") {
			then, _ = time.Parse("02 Jan 06 15:04 +0000", t)
		} else if strings.HasSuffix(t, "GMT") {
			then, _ = time.Parse("02 Jan 06 15:04 GMT", t)
		} else if strings.HasSuffix(t, "UTC") {
			then, _ = time.Parse("02 Jan 06 15:04 UTC", t)
		} else if strings.HasSuffix(t, "CST") {
			then, _ = time.Parse("02 Jan 06 15:04 CST", t)
		} else if strings.HasSuffix(t, "0400") {
			then, _ = time.Parse("02 Jan 06 15:04 -0400", t)
		} else if strings.HasSuffix(t, "Z") {
			then, _ = time.Parse(time.RFC3339, t)
		} else if strings.HasSuffix(t, "0800") {
			then, _ = time.Parse("02 Jan 06 15:04 +0800", t)
		}
	}
	return then
}
