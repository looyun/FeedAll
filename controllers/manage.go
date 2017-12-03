package controllers

import (
	"github.com/looyun/feedall/models"
	"fmt"
	"github.com/mmcdole/gofeed"
	"gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

func AddFeed(c *macaron.Context) bool {
	if !CheckLogin(c) {
		c.HTML(200, "login")
		return false
	}
	username, _ := c.GetSecureCookie("username")
	feedurl := StandarFeed(c.Query("feedurl"))
	feedlists := make([]*models.FeedList, 0)
	fmt.Println("start judge!")
	//Judge if feed existed in feedlist.
	if !models.GetFeedList(models.FeedLists, bson.M{"feedLink": feedurl}, &feedlists) {
		models.GetFeedList(models.FeedLists, bson.M{"feedLink": Prewww(feedurl)}, &feedlists)
		feedurl = Prewww(feedurl)
	}
	if len(feedlists) != 0 {
		fmt.Println("feeds existed!")

		//Judge if feed existed in user's feedurl.
		if models.GetUserInfo(models.Users,
			bson.M{"username": username, "feedLink": feedurl},
			nil) {
			fmt.Println("feeds existed in user's feedlist!")
			return true
		} else {
			return models.UpdateUserFeed(models.Users,
				bson.M{"username": username},
				bson.M{"$push": bson.M{"feedLink": feedurl}})
		}
	} else {
		fmt.Println("Parse feeds!")
		fb := gofeed.NewParser()
		fmt.Println(feedurl)
		value, err := fb.ParseURL(feedurl)
		if err != nil {
			fmt.Println("Parse err: ", err)
			return false
		}
		data, err := bson.Marshal(value)
		if err != nil {
			fmt.Println(err)
			return false
		}
		feed := models.Feed{}
		err = bson.Unmarshal(data, &feed)
		if err != nil {
			fmt.Println(err)
			return false
		}

		feed.FeedLink = feedurl

		for _, v := range feed.Items {
			if v.Content == "" {

				if v.Extensions != nil && v.Extensions["content"] != nil {
					v.Content = v.Extensions["content"]["encoded"][0].Value
				} else {
					v.Content = v.Description
				}
			}
			v.Content = DecodeImg(v.Content, feed.Link)
			if v.Published == "" {
				v.Published = v.Updated
			}
			publishedParsed := ParseDate(v.Published)
			v.PublishedParsed = strconv.FormatInt(publishedParsed.Unix(), 10)
		}
		models.Insert(models.Feeds, feed)
		fmt.Println("inserted feeds!")

		return models.Insert(models.FeedLists, bson.M{"feedLink": feedurl}) &&
			models.UpdateUserFeed(models.Users,
				bson.M{"username": username},
				bson.M{"$push": bson.M{"link": feed.Link}}) &&
			models.UpdateUserFeed(models.Users,
				bson.M{"username": username},
				bson.M{"$push": bson.M{"feedLink": feedurl}})

	}

}

func DelFeed(c *macaron.Context) bool {
	if !CheckLogin(c) {
		c.HTML(200, "login")
		return false
	}
	username, _ := c.GetSecureCookie("username")
	if c.Query("feedurl") != "" {
		feedurl := StandarFeed(c.Query("feedurl"))
		return models.UpdateUserFeed(models.Users,
			bson.M{"username": username},
			bson.M{"$pull": bson.M{"link": feedurl}})
	} else {
		fmt.Println("Feedurl can't be blank!")
		return false
	}
}

func Prewww(s string) string {
	if strings.HasPrefix(s, "http://") {
		s = s[7:]

		if strings.HasPrefix(s, "www") {
			s = "http://" + s[3:]
		} else {
			s = "http://www" + s
		}
	}

	if strings.HasPrefix(s, "https://") {
		s = s[8:]

		if strings.HasPrefix(s, "www") {
			s = "https://" + s[3:]
		} else {
			s = "https://www" + s
		}
	}
	return s
}

func StandarFeed(s string) string {
	if strings.HasSuffix(s, "/") {
		l := len(s)
		s = s[:l-1]
	}
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "http://" + s
	}
	return s
}

func DecodeEntities(str string) string {
	str = strings.Replace(str, "&lt;", "<", -1)
	str = strings.Replace(str, "&gt;", ">", -1)
	str = strings.Replace(str, "&quot;", "\"", -1)
	str = strings.Replace(str, "&apos;", "'", -1)
	str = strings.Replace(str, "&amp;", "&", -1)
	return str
}

func DecodeImg(str string, link string) string {
	str = strings.Replace(str, "&#34;", "\"", -1)
	str = strings.Replace(str, "&quot;", "\"", -1)
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
