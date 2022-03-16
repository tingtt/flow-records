package main

import (
	"flag"
	"flow-records/jwt"
	"flow-records/mysql"
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func getIntEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		var intValue, err = strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Priority: command line params > env variables > default value
var (
	port        = flag.Int("port", getIntEnv("PORT", 1323), "Server port")
	logLevel    = flag.Int("log-level", getIntEnv("LOG_LEVEL", 2), "Log level (1: 'DEBUG', 2: 'INFO', 3: 'WARN', 4: 'ERROR', 5: 'OFF', 6: 'PANIC', 7: 'FATAL'")
	gzipLevel   = flag.Int("gzip-level", getIntEnv("GZIP_LEVEL", 6), "Gzip compression level")
	mysqlHost   = flag.String("mysql-host", getEnv("MYSQL_HOST", "db"), "MySQL host")
	mysqlPort   = flag.Int("mysql-port", getIntEnv("MYSQL_PORT", 3306), "MySQL port")
	mysqlDB     = flag.String("mysql-database", getEnv("MYSQL_DATABASE", "flow-records"), "MySQL database")
	mysqlUser   = flag.String("mysql-user", getEnv("MYSQL_USER", "flow-records"), "MySQL user")
	mysqlPasswd = flag.String("mysql-password", getEnv("MYSQL_PASSWORD", ""), "MySQL password")
	jwtIssuer   = flag.String("jwt-issuer", getEnv("JWT_ISSUER", "flow-users"), "JWT issuer")
	jwtSecret   = flag.String("jwt-secret", getEnv("JWT_SECRET", ""), "JWT secret")
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: *gzipLevel,
	}))
	e.Logger.SetLevel(log.Lvl(*logLevel))
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.JwtCustumClaims{},
		SigningKey: []byte(*jwtSecret),
	}))

	// Setup db client instance
	e.Logger.Info(mysql.SetDSNTCP(*mysqlUser, *mysqlPasswd, *mysqlHost, *mysqlPort, *mysqlDB))

	// Restricted routes
	e.POST("/", post)
	e.PUT("/", putList)
	e.GET("/", getList)
	e.PATCH("/:id", patch)
	e.DELETE("/:id", delete)
	e.POST("/changelogs", changeLogPost)
	e.GET("/changelogs", changeLogGetList)
	e.GET("/changelogs/:id", changeLogGet)
	e.PATCH("/changelogs/:id", changeLogPatch)
	e.DELETE("/changelogs/:id", changeLogDelete)
	e.POST("/schemes", schemePost)
	e.GET("/schemes", schemeGetList)
	e.GET("/schemes/:id", schemeGet)
	e.PATCH("/schemes/:id", schemePatch)
	e.DELETE("/schemes/:id", schemeDelete)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
