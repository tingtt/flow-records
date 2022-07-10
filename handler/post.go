package handler

import (
	"encoding/json"
	"flow-records/flags"
	"flow-records/jwt"
	"flow-records/record"
	"flow-records/scheme"
	"flow-records/utils"
	"fmt"
	"net/http"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Post(c echo.Context) error {
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
	postTmp := new(interface{})
	if err = c.Bind(postTmp); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Marshal post body
	b, err := json.Marshal(postTmp)
	if err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	postMultiple := new(record.MultiplePostBody)
	err = json.Unmarshal(b, postMultiple)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request body
	err = c.Validate(postMultiple)
	if err == nil {
		// Multiple records

		// Check todo id
		if postMultiple.TodoId != nil {
			status, err := utils.HttpGet(fmt.Sprintf("%s/%d", *flags.Get().ServiceUrlTodos, *postMultiple.TodoId), &u.Raw)
			if err != nil {
				// 500: Internal server error
				c.Logger().Error(err)
				return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
			}
			if status != http.StatusOK {
				// 400: Bad request
				c.Logger().Debugf("todo id: %d does not exist", *postMultiple.TodoId)
				return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("todo id: %d does not exist", *postMultiple.TodoId)}, "	")
			}
		}

		// Check schemeIds
		for _, records := range postMultiple.Records {
			_, notFound, err := scheme.Get(userId, records.SchemeId, scheme.GetQuery{})
			if err != nil {
				// 500: Internal server error
				c.Logger().Debug(err)
				return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
			}
			if notFound {
				// 400: Bad request
				c.Logger().Debugf("scheme id: %d does not exists", records.SchemeId)
				return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("scheme id: %d does not exists", records.SchemeId)}, "	")
			}
		}

		records, err := record.PostMultiple(userId, *postMultiple)
		if err != nil {
			// 500: Internal server error
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// 200: Success
		if records == nil {
			return c.JSONPretty(http.StatusOK, []interface{}{}, "	")
		}
		return c.JSONPretty(http.StatusOK, records, "	")

	}

	// Single record

	post := new(record.PostBody)
	err = json.Unmarshal(b, post)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request body
	if err = c.Validate(post); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// Check todo id
	status, err := utils.HttpGet(fmt.Sprintf("%s/%d", *flags.Get().ServiceUrlTodos, post.TodoId), &u.Raw)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if status != http.StatusOK {
		// 400: Bad request
		c.Logger().Debugf("todo id: %d does not exist", post.TodoId)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("todo id: %d does not exist", post.TodoId)}, "	")
	}

	// Check schemeId
	_, notFound, err := scheme.Get(userId, post.SchemeId, scheme.GetQuery{})
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 400: Bad request
		c.Logger().Debugf("scheme id: %d does not exists", post.SchemeId)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("scheme id: %d does not exists", post.SchemeId)}, "	")
	}

	r, err := record.Post(userId, *post)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, r, "	")
}
