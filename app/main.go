package main

import (
	"crud-product/config"
	"crud-product/delivery/rest"
	"crud-product/repository"
	"crud-product/usecase"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
)


// @title Echo Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	// Init config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Init echo framework
	e := echo.New()

	// Init DB
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)

	db, err := sql.Open("mysql", mysqlInfo)
	if err != nil {
		fmt.Println("Error connect to db")
		log.Fatal(err)
	}

	defer db.Close()

	// Init repository
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Init usecase
	productUsecae := usecase.NewProduct(productRepo)
	userUsecase := usecase.NewUser(userRepo)

	// Init handler
	rest.NewHandler(e, productUsecae, userUsecase)

	e.GET("/", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}