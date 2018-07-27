package main

import (
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"time"

	"github.com/go-macaron/gzip"
	"github.com/looyun/feedall/controllers"
	"github.com/looyun/feedall/middleware"
	"github.com/looyun/feedall/models"
	"github.com/looyun/feedall/parse"
	macaron "gopkg.in/macaron.v1"
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
	m.Use(gzip.Gziper())
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
	m.Post("/login", func(ctx *macaron.Context) {
		token, err := controllers.Login(ctx)
		if err != nil {
			fmt.Println(err)
			ctx.Error(400, "error")
			return
		}
		ctx.JSON(200, map[string]string{"token": token})
	})
	m.Post("/signup", func(ctx *macaron.Context) {
		err := controllers.Signup(ctx)
		if err != nil {
			fmt.Println(err)
			ctx.Error(400, "error")
		} else {
			ctx.Status(200)
		}
	})
	m.Post("/item/:feedlink", func(ctx *macaron.Context) {
	})

	m.Group("/api", func() {
		m.Group("/my", func() {

			m.Get("/feeds", func(ctx *macaron.Context) {
				feeds, err := controllers.GetUserFeeds(ctx)
				if err != nil {
					fmt.Println(err)
					ctx.Error(400, "error")
				} else {
					ctx.JSON(200, &feeds)
				}

			})
			m.Get("/items", func(ctx *macaron.Context) {
				items, err := controllers.GetUserItems(ctx)
				if err != nil {
					fmt.Println(err)
					ctx.Error(400, "error")
				} else {
					ctx.JSON(200, &items)
				}

			})
			m.Get("/stars", func(ctx *macaron.Context) {
				items, err := controllers.GetStarItems(ctx)
				if err != nil {
					fmt.Println(err)
					ctx.Error(400, "error")
				} else {
					ctx.JSON(200, &items)
				}

			})
			m.Post("/subscribe", func(ctx *macaron.Context) {
				err := controllers.Subscribe(ctx)
				if err != nil {
					fmt.Println(err)
					ctx.Error(400, "error")
				} else {
					ctx.Status(200)
				}
			})

			// m.Post("/del", func(ctx *macaron.Context) {
			// 	if controllers.DelFeed(ctx) {
			// 		fmt.Println("Delete feed succeed!")
			// 		ctx.Redirect("/manage")
			// 	} else {
			// 		fmt.Println("Delete feed false!")
			// 		ctx.Redirect("/manage")
			// 	}
			// })
		}, middleware.ValidateJWTToken())

		m.Get("/feeds/recommand/:n:int", func(ctx *macaron.Context) {
			feeds := controllers.GetFeeds(ctx)
			ctx.JSON(200, &feeds)

		})

		m.Get("/feed/:id/", func(ctx *macaron.Context) {
			feed := controllers.GetFeed(ctx)
			ctx.JSON(200, &feed)

		})
		m.Get("/feed/:id/items", func(ctx *macaron.Context) {
			items := controllers.GetFeedItems(ctx)
			ctx.JSON(200, &items)

		})
		m.Get("/item/:id", func(ctx *macaron.Context) {
			item := controllers.GetItem(ctx)
			ctx.JSON(200, &item)
		})
		m.Get("/items/random/:n:int", func(ctx *macaron.Context) {
			numbers := ctx.ParamsInt(":n")
			items := controllers.GetRandomItem(ctx, numbers)
			ctx.JSON(200, &items)
		})
		m.Get("/items/recommand/:n:int", func(ctx *macaron.Context) {
			items := controllers.GetItems(ctx, 5)
			ctx.JSON(200, &items)
		})

	})
	m.Run()

}
