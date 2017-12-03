package main

import (
	"github.com/looyun/feedall/controllers"
	"github.com/looyun/feedall/models"
	"github.com/looyun/feedall/parse"
	"fmt"
	"gopkg.in/macaron.v1"
	"html/template"
	"net/url"
	"strconv"
	"time"
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
	}))
	m.SetDefaultCookieSecret("feedall")
	m.Get("/", func(ctx *macaron.Context) {
		if !controllers.CheckLogin(ctx) {
			ctx.HTML(200, "welcome")
			return
		}
		ctx.Data["IsLogin"] = controllers.CheckLogin(ctx)
		controllers.GetUserFeed(ctx)

		ctx.HTML(200, "index")
	})

	m.Get("/feed", func(ctx *macaron.Context) {
		if !controllers.CheckLogin(ctx) {
			ctx.HTML(200, "welcome")
			return
		}
		ctx.Data["IsLogin"] = controllers.CheckLogin(ctx)
		controllers.GetUserFeed(ctx)
		ctx.HTML(200, "index")
	})
	m.Post("/feed", func(ctx *macaron.Context) {
		if !controllers.CheckLogin(ctx) {
			ctx.HTML(200, "welcome")
			return
		}
		ctx.Data["IsLogin"] = controllers.CheckLogin(ctx)
		controllers.GetUserFeed(ctx)
	})

	m.Post("/item/*", func(ctx *macaron.Context) {
		if !controllers.CheckLogin(ctx) {
			ctx.HTML(200, "welcome")
			return
		}
		ctx.Data["IsLogin"] = controllers.CheckLogin(ctx)
		controllers.GetItemContent(ctx)
	})

	m.Get("/item/*", func(ctx *macaron.Context) {
		if !controllers.CheckLogin(ctx) {
			ctx.HTML(200, "welcome")
			return
		}
		ctx.Data["IsLogin"] = controllers.CheckLogin(ctx)
		controllers.GetItemContent(ctx)
		ctx.HTML(200, "index")
	})

	m.Post("/add", func(ctx *macaron.Context) {
		if !controllers.CheckLogin(ctx) {
			ctx.HTML(200, "login")
			return
		}
		if controllers.AddFeed(ctx) {
			ctx.Redirect("/")
		} else {
			fmt.Println("Add feed false!")
			ctx.Redirect("/")
		}
	})
	m.Get("/manage", func(ctx *macaron.Context) {
		if !controllers.CheckLogin(ctx) {
			ctx.Redirect("/user/login")
			return
		}
		ctx.Data["IsLogin"] = controllers.CheckLogin(ctx)
		controllers.GetUserFeed(ctx)
		ctx.HTML(200, "manage")
	})

	m.Post("/del", func(ctx *macaron.Context) {
		if !controllers.CheckLogin(ctx) {
			ctx.HTML(200, "login")
			return
		}
		if controllers.DelFeed(ctx) {
			fmt.Println("Delete feed succeed!")
			ctx.Redirect("/manage")
		} else {
			fmt.Println("Delete feed false!")
			ctx.Redirect("/manage")
		}
	})
	m.Group("/user", func() {

		m.Get("/login", func(ctx *macaron.Context) {
			if ctx.Query("exit") == "true" {
				ctx.SetCookie("username", "")

				ctx.SetCookie("password", "")
				ctx.Redirect("/")
				return
			}
			ctx.Data["IsLogin"] = controllers.CheckLogin(ctx)
			ctx.HTML(200, "login")
		})
		m.Post("/login", func(ctx *macaron.Context) {
			if controllers.Login(ctx) {
				ctx.Redirect("/")
				return
			} else {
				ctx.HTML(200, "login")
				return
			}
		})
		m.Get("/test", func(ctx *macaron.Context) {

		})
		m.Get("/register", func(ctx *macaron.Context) {
			ctx.Data["IsLogin"] = controllers.CheckLogin(ctx)
			ctx.HTML(200, "register")
		})
		m.Post("/register", func(ctx *macaron.Context) {
			controllers.Register(ctx)
		})
	})
	m.Run()

}
