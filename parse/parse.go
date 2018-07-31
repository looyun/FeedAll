package parse

import (
	"fmt"
	"time"

	"github.com/looyun/feedall/controllers"
	"github.com/looyun/feedall/models"
)

func Parse() {
	for {
		timer := time.NewTimer(600 * time.Second)
		fmt.Println("start parse!")
		feeds := make([]models.Feed, 0)
		err := models.FindAll(models.Feeds, nil, &feeds)
		if err != nil {
			fmt.Println(<-timer.C)
			continue
		} else {
			Finish := make(chan string)
			for _, feed := range feeds {
				go func(feed models.Feed) {
					err := controllers.UpdateItems(feed)
					if err != nil {
						fmt.Println(err)
					}
					Finish <- feed.FeedLink
				}(feed)
			}
			for _, _ = range feeds {
				fmt.Println(<-Finish)
			}
		}
		fmt.Println("OK!")
		fmt.Println(<-timer.C)
	}
}
