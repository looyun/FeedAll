package main

import (
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"time"

	"github.com/looyun/feedall/controllers"
	"github.com/looyun/feedall/models"
	"github.com/looyun/feedall/parse"
	macaron "gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2/bson"
)

const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func main() {

	models.Init()
	m := macaron.Classic()
	go parse.Parse()

	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs: []template.FuncMap{map[string]interface{}{
			"str2html": func(raw string) template.HTML {
				return template.HTML(raw)

			},
			"UrlParse": func(raw string) string {
				return url.QueryEscape(raw)

			},
			"TimeSince": func(s string) string {
				now := time.Now()
				i, _ := strconv.ParseInt(s, 10, 64)
				then := time.Unix(i, 0)
				diff := now.Unix() - then.Unix()
				if then.After(now) {
					diff = then.Unix() - now.Unix()
				}
				switch {
				case diff <= 0:
					return "now"
				case diff <= 2:
					return "1s"
				case diff < 1*Minute:
					return strconv.FormatInt(diff, 10) + "s"

				case diff < 2*Minute:
					return "1m"
				case diff < 1*Hour:
					return strconv.FormatInt(diff/Minute, 10) + "m"

				case diff < 2*Hour:
					return "1h"
				case diff < 1*Day:
					return strconv.FormatInt(diff/Hour, 10) + "h"

				case diff < 2*Day:
					return "1d"
				case diff < 1*Week:
					return strconv.FormatInt(diff/Day, 10) + "d"

				case diff < 2*Week:
					return "1w"
				default:
					return then.Month().String()[:3] + " " + strconv.Itoa(then.Year())
				}
			},
		}},
		IndentJSON: true,
	}))
	m.SetDefaultCookieSecret("feedall")
	m.Group("/api", func() {
		m.Group("/user", func() {

			m.Post("/login", func(ctx *macaron.Context) {
				if controllers.Login(ctx) {
					ctx.Redirect("/")
				} else {
					ctx.JSON(200, "login")
				}
				return
			})
			m.Post("/register", func(ctx *macaron.Context) {
				controllers.Register(ctx)
			})

			m.Post("/feed", func(ctx *macaron.Context) {
				ctx.Data["IsLogin"] = controllers.CheckLogin(ctx)
				controllers.GetUserFeed(ctx)
			})

			m.Post("/add", func(ctx *macaron.Context) {
				controllers.AddFeed(ctx)
			})

			m.Post("/del", func(ctx *macaron.Context) {
				if controllers.DelFeed(ctx) {
					fmt.Println("Delete feed succeed!")
					ctx.Redirect("/manage")
				} else {
					fmt.Println("Delete feed false!")
					ctx.Redirect("/manage")
				}
			})
		})

		m.Get("/feeds", func(ctx *macaron.Context) {
			feeds := controllers.GetFeeds(ctx)
			ctx.JSON(200, &feeds)

		})
		// m.Get("/feeds/hot/:n:int", func(ctx *macaron.Context) {
		// 	feed := controllers.GetFeed(ctx)
		// 	ctx.JSON(200, &feed)

		// })

		m.Get("/feed/:feedlink/", func(ctx *macaron.Context) {
			feed := controllers.GetFeed(ctx)
			ctx.JSON(200, &feed)

		})
		m.Get("/feed/:feedlink/items", func(ctx *macaron.Context) {
			page := ctx.QueryInt("page")
			if page == 0 {
				page = 1
			}
			per_page := ctx.QueryInt("per_page")
			if per_page == 0 {
				per_page = 20
			}

			feedlink := controllers.ParseURL(ctx.Params(":feedlink"))
			feed := models.Feed{}
			models.FindOne(models.Feeds,
				bson.M{"feedLink": feedlink},
				&feed)

			items := []bson.M{}
			models.Items.Find(bson.M{"feedID": feed.ID}).Sort("-publishedParsed").Skip(per_page * (page - 1)).Limit(per_page).All(&items)

			ctx.JSON(200, &items)

		})
		m.Get("/item/:itemlink", func(ctx *macaron.Context) {
			item := controllers.GetItem(ctx)
			ctx.JSON(200, &item)
		})
		m.Get("/items/random/:n:int", func(ctx *macaron.Context) {
			numbers := ctx.ParamsInt(":n")
			items := controllers.GetRandomItem(ctx, numbers)
			ctx.JSON(200, &items)
		})
		m.Get("/items/new/:n:int", func(ctx *macaron.Context) {
			numbers := ctx.ParamsInt(":n")
			items := controllers.GetLatestItem(ctx, numbers)
			ctx.JSON(200, &items)
		})
		// m.Get("/items/hot", func(ctx *macaron.Context) {
		// 	items := controllers.GetItemSample(ctx, 5)
		// 	ctx.JSON(200, &items)
		// })

	})
	m.Run()

}
