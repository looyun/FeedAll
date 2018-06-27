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
					feed, err := fb.ParseURL(u.FeedLink)
					if err != nil {
						fmt.Println("Parse err: ", err)
						Finish <- u.FeedLink
					} else {

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

							var item_ids []int
							item_id := models.InsertItem(models.Items,
								bson.M{"link": v.Link},
								v)
							if item_id != nil {
								item_ids = append(item_ids, item_id.(int))
							}
						}
						feed.ItemIDs = item_ids
						models.UpdateFeed(models.Feeds,
							bson.M{"feedlink": feed.FeedLink},
							feed)
						fmt.Println("updatefeed OK!")
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
