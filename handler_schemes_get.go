package main

import (
	"flow-records/jwt"
	"flow-records/scheme"
	"net/http"
	"strconv"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type GetQuery struct {
	Embed *string `query:"embed" validate:"omitempty,oneof=records record.changelog"`
	Start *string `query:"start" validate:"omitempty,datetime"`
	End   *string `query:"end" validate:"omitempty,datetime"`
}

func schemeGet(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// id
	idStr := c.Param("id")
	// string -> uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 404: Not found
		return echo.ErrNotFound
	}

	// Bind request query
	query := new(GetQuery)
	if err = c.Bind(query); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request query
	if err = c.Validate(query); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}
	var start, end *time.Time
	if query.Start != nil {
		startTmp, err := datetimeStrConv(*query.Start)
		if err != nil {
			// 400: Bad request
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}
		start = &startTmp
	}
	if query.End != nil {
		endTmp, err := datetimeStrConv(*query.End)
		if err != nil {
			// 400: Bad request
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}
		end = &endTmp
	}
	queryParsed := scheme.GetQuery{Embed: query.Embed, Start: start, End: end}

	// Read db
	s, notFound, err := scheme.Get(userId, id, queryParsed)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("scheme not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "scheme not found"}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, s, "	")
}
