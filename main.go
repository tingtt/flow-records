package main

import (
	"flow-records/flags"
	"flow-records/handler"
	"flow-records/jwt"
	"flow-records/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

type CustomValidator struct {
	validator *validator.Validate
}

func DatetimeStrValidation(fl validator.FieldLevel) bool {
	_, err1 := time.Parse("2006-1-2T15:4:5", fl.Field().String())
	_, err2 := time.Parse(time.RFC3339, fl.Field().String())
	_, err3 := strconv.ParseUint(fl.Field().String(), 10, 64)
	return err1 == nil || err2 == nil || err3 == nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	cv.validator.RegisterValidation("datetime", DatetimeStrValidation)

	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

func main() {
	// Get command line params / env variables
	f := flags.Get()

	//
	// Setup echo and middlewares
	//

	// Echo instance
	e := echo.New()

	// Gzip
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: int(*f.GzipLevel),
	}))

	// Log level
	e.Logger.SetLevel(log.Lvl(*f.LogLevel))

	// Validator instance
	e.Validator = &CustomValidator{validator: validator.New()}

	// JWT
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.JwtCustumClaims{},
		SigningKey: []byte(*f.JwtSecret),
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/-/readiness"
		},
	}))

	//
	// Check health of external service
	//

	// flow-projects
	if *flags.Get().ServiceUrlProjects == "" {
		e.Logger.Fatal("`--service-url-projects` option is required")
	}
	if status, err := utils.HttpGet(*flags.Get().ServiceUrlProjects+"/-/readiness", nil); err != nil {
		e.Logger.Fatalf("failed to check health of external service `flow-projects` %s", err)
	} else if status != http.StatusOK {
		e.Logger.Fatal("failed to check health of external service `flow-projects`")
	}
	// flow-todos
	if *flags.Get().ServiceUrlTodos == "" {
		e.Logger.Fatal("`--service-url-todos` option is required")
	}
	if status, err := utils.HttpGet(*flags.Get().ServiceUrlTodos+"/-/readiness", nil); err != nil {
		e.Logger.Fatalf("failed to check health of external service `flow-todos` %s", err)
	} else if status != http.StatusOK {
		e.Logger.Fatal("failed to check health of external service `flow-todos`")
	}

	//
	// Routes
	//

	// Health check route
	e.GET("/-/readiness", func(c echo.Context) error {
		return c.String(http.StatusOK, "flow-records is Healthy.\n")
	})

	// Restricted routes
	e.POST("/", handler.Post)
	e.PUT("/", handler.PutList)
	e.GET("/", handler.GetList)
	e.GET("/:id", handler.Get)
	e.PATCH("/:id", handler.Patch)
	e.DELETE("/:id", handler.Delete)
	e.DELETE("/", handler.DeleteAll)
	e.POST("/changelogs", handler.ChangeLogPost)
	e.GET("/changelogs", handler.ChangeLogGetList)
	e.GET("/changelogs/:id", handler.ChangeLogGet)
	e.PATCH("/changelogs/:id", handler.ChangeLogPatch)
	e.DELETE("/changelogs/:id", handler.ChangeLogDelete)
	e.DELETE("/changelogs", handler.ChangeLogDeleteAll)
	e.POST("/schemes", handler.SchemePost)
	e.GET("/schemes", handler.SchemeGetList)
	e.GET("/schemes/:id", handler.SchemeGet)
	e.PATCH("/schemes/:id", handler.SchemePatch)
	e.DELETE("/schemes/:id", handler.SchemeDelete)
	e.DELETE("/schemes", handler.SchemeDeleteAll)

	//
	// Start echo
	//
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *f.Port)))
}
