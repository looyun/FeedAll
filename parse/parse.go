package parse

import (
	"fmt"
	"strconv"
	"time"

	"github.com/looyun/feedall/controllers"
	"github.com/looyun/feedall/models"
	"github.com/mmcdole/gofeed"
	"gopkg.in/mgo.v2/bson"
)

func Parse() {
	for {
		timer := time.NewTimer(600 * time.Second)
		fmt.Println("start parse!")
		feeds := make([]*models.Feed, 0)
		err := models.FindAll(models.Feeds, nil, &feeds)
		if err != nil {
			fmt.Println(<-timer.C)
			continue
		} else {
			Finish := make(chan string)
			fb := gofeed.NewParser()
			for _, u := range feeds {
				go func(u *models.Feed) {
					origin_feed, err := fb.ParseURL(u.FeedLink)
					if err != nil {
						fmt.Println("Parse err: ", err)
						Finish <- u.FeedLink
					} else {
						data, err := bson.Marshal(origin_feed)
						if err != nil {
							fmt.Println(err)
						}
						var items struct {
							Items []*models.Item `bson:"items"`
						}
						err = bson.Unmarshal(data, &items)
						if err != nil {
							fmt.Println(err)
						}
						feed := models.Feed{}
						models.FindOne(models.Feeds,
							bson.M{"feedLink": u.FeedLink},
							&feed)
						for _, v := range items.Items {
							if v.Content == "" {
								if v.Extensions != nil && v.Extensions["content"] != nil {
									v.Content = v.Extensions["content"]["encoded"][0].Value
								} else {
									v.Content = v.Description
								}
							}
							v.Content = controllers.DecodeImg(v.Content, u.Link)
							if v.Published == "" {
								v.Published = v.Updated
							}
							publishedParsed := controllers.ParseDate(v.Published)
							v.PublishedParsed = strconv.FormatInt(publishedParsed.Unix(), 10)
							v.FeedID = feed.ID
							v.ID = bson.NewObjectId()

							info, err := models.Upsert(models.Items,
								bson.M{"link": v.Link},
								v)
							if err != nil {
								fmt.Println(err)
							}
							if info == nil {
								continue
							}
						}
						Finish <- u.FeedLink
					}
				}(u)
			}
			for _, _ = range feeds {
				fmt.Println(<-Finish)
			}
		}
		fmt.Println("OK!")
		fmt.Println(<-timer.C)
	}
}
