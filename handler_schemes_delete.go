package main

import (
	"flow-records/jwt"
	"flow-records/scheme"
	"net/http"
	"strconv"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func schemeDelete(c echo.Context) error {
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

	notFound, err := scheme.Delete(userId, id)
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

	// 204: No content
	return c.JSONPretty(http.StatusNoContent, map[string]string{"message": "Deleted"}, "	")
}
