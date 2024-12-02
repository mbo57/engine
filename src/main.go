package main

import (
	"app/handler"
	"app/infra/disk"
	"app/router"
	"app/usecase"
	"app/util/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	logger := logger.New()
	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	h := handler.NewIndexHandler(
		logger,
		usecase.NewIndexUsecase(
			logger,
			disk.NewIndexRepository("data"),
		),
	)
	router.NewRouter(e, h)
	e.Logger.Fatal(e.Start(":1323"))
}
