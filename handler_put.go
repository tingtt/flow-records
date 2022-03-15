package main

import (
	"flow-records/jwt"
	"flow-records/record"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type putQuery struct {
	TodoId uint64 `query:"todo_id" validate:"required,gte=1"`
}

func putList(c echo.Context) error {
	// Check `Content-Type`
	if c.Request().Header.Get("Content-Type") != "application/json" {
		// 415: Invalid `Content-Type`
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported media type"}, "	")
	}

	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	records, err := record.Put(userId, putQuery{}.TodoId, record.PutBody{})
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, records, "	")
}