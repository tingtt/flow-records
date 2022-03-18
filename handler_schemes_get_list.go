package main

import (
	"flow-records/jwt"
	"flow-records/scheme"
	"net/http"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type GetListQuery struct {
	ProjectId *uint64 `query:"project_id" validate:"omitempty,gte=1"`
	Start     *string `query:"start" validate:"omitempty,datetime"`
	End       *string `query:"end" validate:"omitempty,datetime"`
	Embed     *string `query:"embed" validate:"omitempty,oneof=records record.changelog"`
}

func schemeGetList(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// Bind request query
	query := new(GetListQuery)
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
	queryParsed := scheme.GetListQuery{ProjectId: query.ProjectId, Embed: query.Embed, Start: start, End: end}

	schemes, err := scheme.GetList(userId, queryParsed)
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	if schemes == nil {
		return c.JSONPretty(http.StatusOK, []interface{}{}, "	")
	}
	return c.JSONPretty(http.StatusOK, schemes, "	")
}
