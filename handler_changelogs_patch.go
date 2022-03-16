package main

import (
	"flow-records/changelog"
	"flow-records/jwt"
	"net/http"
	"strconv"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func changeLogPatch(c echo.Context) error {
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

	// id
	idStr := c.Param("id")
	// string -> uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 404: Not found
		return echo.ErrNotFound
	}

	cl, notFound, err := changelog.Patch(userId, id, changelog.PatchBody{})
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("changelog not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "changelog not found"}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, cl, "	")
}