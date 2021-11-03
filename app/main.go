package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/apm-dev/vending-machine/pkg/logger"
	"github.com/apm-dev/vending-machine/user"
	userPgsql "github.com/apm-dev/vending-machine/user/data/pgsql"
	userRest "github.com/apm-dev/vending-machine/user/presentation/rest"
	"github.com/apm-dev/vending-machine/user/presentation/rest/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	logLevel := logger.INFO
	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
		logLevel = logger.DEBUG
	}
	logger.SetLogger(logger.NewLogcat(logLevel))
}

func main() {
	// infrastructure (db,...)
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fatalOnError(err)
	dbConn, err := db.DB()
	fatalOnError(err)
	err = dbConn.Ping()
	fatalOnError(err)

	defer func() {
		fatalOnError(dbConn.Close())
	}()

	// data (repository)
	err = db.AutoMigrate(&userPgsql.User{}, &userPgsql.JWT{})
	fatalOnError(err)

	ur := userPgsql.InitUserRepository(db)
	jr := userPgsql.InitJwtRepository(db)
	jwt := user.NewJWTManager(
		viper.GetString("jwt.secret"),
		time.Duration(viper.GetInt("jwt.duration"))*time.Second,
	)

	depositTimeout := time.Duration(viper.GetInt("deposit.timeout")) * time.Second

	// services (usecase)
	us := user.InitService(ur, jr, jwt, depositTimeout)

	// presentation (delivery/controller)
	e := echo.New()
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("1M"))
	// echo validator
	e.Validator = &CustomValidator{validator: validator.New()}
	// echo middlewares
	authMiddleware := middlewares.InitUserMiddleware(us)
	ag := e.Group("", authMiddleware.JwtAuth)

	// rest(http) handlers
	userRest.InitUserHandler(e, ag, us)

	log.Fatal(e.Start(viper.GetString("server.address")))
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func fatalOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
