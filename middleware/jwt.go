package middleware

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/looyun/feedall/controllers"
	"github.com/looyun/feedall/models"
	macaron "gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2/bson"
)

func ValidateJWTToken() macaron.Handler {
	return func(ctx *macaron.Context) {
		token, err := request.ParseFromRequest(ctx.Req.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return controllers.TokenSecure, nil
			})

		if err == nil {
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				username := claims["username"].(string)

				user := models.User{}
				err := models.FindOne(models.Users, bson.M{"username": username}, &user)
				if err != nil {
					fmt.Println(err)
					ctx.Status(http.StatusUnauthorized)
				}
				ctx.Data["user"] = user

				ctx.Next()
			} else {
				fmt.Println("Token not valid")
				ctx.Status(http.StatusUnauthorized)
			}
		} else {
			fmt.Println(err)
			ctx.Status(http.StatusUnauthorized)
		}
	}
}
