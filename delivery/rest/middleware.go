package rest

import (
	"net/http"
	"strings"

	"crud-product/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Exception struct {
	Message string `json:"message"`
}

func JwtVerify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		var header = c.Request().Header.Get("x-access-token") // Grab the token from the header

		header = strings.TrimSpace(header)

		if header == "" {
			// Token is missing, returns with error code 403 Unauthorized
			return c.JSON(http.StatusForbidden, Exception{
				Message: "Missing auth token"},
			)
		}

		tk := &model.Token{}

		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			return c.JSON(http.StatusForbidden, Exception{
				Message: err.Error()},
			)
		}

		c.Set("user", tk)

		return next(c)
	}
}