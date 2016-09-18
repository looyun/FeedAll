package parse

import (
	"feedall/controllers"
	"feedall/models"
	"fmt"
	"github.com/mmcdole/gofeed"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

func Parse() {
	for {
		timer := time.NewTimer(60 * time.Second)

		fmt.Println("start parse!")
		feedlist := make([]*models.FeedList, 0)
		if !models.GetFeedList(models.FeedLists, nil, &feedlist) {
			fmt.Println(<-timer.C)
			continue
		} else {
			Finish := make(chan string)

			fb := gofeed.NewParser()
			for _, u := range feedlist {
				go func(u *models.FeedList) {
					value, err := fb.ParseURL(u.FeedLink)
					if err != nil {
						fmt.Println("Parse err: ", err)
						Finish <- u.FeedLink
					} else {
						data, err := bson.Marshal(value)
						if err != nil {
							fmt.Println(err)
						}
						feed := models.Feed{}
						err = bson.Unmarshal(data, &feed)
						if err != nil {
							fmt.Println(err)
						}

						for _, v := range feed.Items {
							if v.Content == "" {
								if v.Extensions != nil && v.Extensions["content"] != nil {
									v.Content = v.Extensions["content"]["encoded"][0].Value
								} else {
									v.Content = v.Description
								}
							}
							v.Content = controllers.DecodeImg(v.Content, feed.Link)
							if v.Published == "" {
								v.Published = v.Updated
							}
							publishedParsed := controllers.ParseDate(v.Published)
							v.PublishedParsed = strconv.FormatInt(publishedParsed.Unix(), 10)
						}

						for _, v := range feed.Items {
							if !models.GetItem(models.Feeds,
								[]bson.M{
									bson.M{"$match": bson.M{"items.link": v.Link}},
								},
								bson.M{}) {
								models.UpdateFeed(models.Feeds,
									bson.M{"title": feed.Title},
									bson.M{"$push": bson.M{"items": v}})
								fmt.Println("updatefeed OK!")
							} else {
								break
							}

						}
						Finish <- u.FeedLink
					}

				}(u)
			}
			for _, _ = range feedlist {
				fmt.Println(<-Finish)
			}
		}
		fmt.Println("OK!")
		fmt.Println(<-timer.C)
	}
}
