package handler

import (
	"flow-records/flags"
	"flow-records/jwt"
	"flow-records/scheme"
	"flow-records/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func SchemePatch(c echo.Context) error {
	// Check `Content-Type`
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		// 415: Invalid `Content-Type`
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported media type"}, "	")
	}

	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
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

	// Check project id
	if post.ProjectId.UInt64 != nil && *post.ProjectId.UInt64 != nil {
		status, err := utils.HttpGet(fmt.Sprintf("%s/%d", *flags.Get().ServiceUrlProjects, **post.ProjectId.UInt64), &u.Raw)
		if err != nil {
			// 500: Internal server error
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if status != http.StatusOK {
			// 400: Bad request
			c.Logger().Debugf("project id: %d does not exist", **post.ProjectId.UInt64)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("project id: %d does not exist", **post.ProjectId.UInt64)}, "	")
		}
	}

	s, notFound, err := scheme.Patch(userId, id, *post)
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
