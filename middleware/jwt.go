package middleware

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/looyun/feedall/controllers"
	macaron "gopkg.in/macaron.v1"
)

func ValidateJWTToken() macaron.Handler {
	fmt.Println("valid")
	return func(ctx *macaron.Context) {
		fmt.Println("func")
		token, err := request.ParseFromRequest(ctx.Req.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return controllers.TokenSecure, nil
			})

		if err == nil {
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				username := claims["username"].(string)
				ctx.Data["username"] = username
				ctx.Next()
			} else {
				ctx.Status(http.StatusUnauthorized)
			}
		} else {
			ctx.Status(http.StatusUnauthorized)
		}
	}
}
