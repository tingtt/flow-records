package handler

import (
	"flow-records/flags"
	"flow-records/jwt"
	"flow-records/scheme"
	"flow-records/utils"
	"fmt"
	"net/http"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func SchemePost(c echo.Context) error {
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

	// Bind request body
	post := new(scheme.PostBody)
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
	if post.ProjectId != nil {
		status, err := utils.HttpGet(fmt.Sprintf("%s/%d", *flags.Get().ServiceUrlProjects, *post.ProjectId), &u.Raw)
		if err != nil {
			// 500: Internal server error
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if status != http.StatusOK {
			// 400: Bad request
			c.Logger().Debugf("project id: %d does not exist", *post.ProjectId)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("project id: %d does not exist", *post.ProjectId)}, "	")
		}
	}

	// Write to db
	s, err := scheme.Post(userId, *post)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, s, "	")
}
