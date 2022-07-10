package handler

import (
	"flow-records/flags"
	"flow-records/jwt"
	"flow-records/record"
	"flow-records/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type putQuery struct {
	TodoId uint64 `query:"todo_id" validate:"required,gte=1"`
}

func PutList(c echo.Context) error {
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

	// Bind and validater request query
	todoId, err := strconv.ParseUint(c.QueryParam("todo_id"), 10, 64)
	if err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}
	query := &putQuery{todoId}
	if err = c.Validate(query); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// Bind request body
	put := new(record.PutBody)
	if err = c.Bind(put); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request body
	type putBodyWrap struct {
		Records record.PutBody `validate:"required,gte=1,dive"`
	}
	if err = c.Validate(putBodyWrap{*put}); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// Check todo id
	status, err := utils.HttpGet(fmt.Sprintf("%s/%d", *flags.Get().ServiceUrlTodos, todoId), &u.Raw)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if status != http.StatusOK {
		// 400: Bad request
		c.Logger().Debugf("todo id: %d does not exist", todoId)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("todo id: %d does not exist", todoId)}, "	")
	}

	records, err := record.Put(userId, query.TodoId, *put)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, records, "	")
}
