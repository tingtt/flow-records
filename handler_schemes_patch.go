package main

import (
	"flow-records/jwt"
	"flow-records/scheme"
	"net/http"
	"strconv"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func schemePatch(c echo.Context) error {
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

	// Bind request body
	post := new(scheme.PatchBody)
	if err = c.Bind(post); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request body
	if err = c.Validate(post); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// TODO: Check projectId

	s, notFound, err := scheme.Patch(userId, id, *post)
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
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
