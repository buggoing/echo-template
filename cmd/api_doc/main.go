package main

import (
	"fmt"

	_ "github.com/PPIO/pi-cloud-monitor-backend/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	docServerPort = 20002
)

func main() {
	e := echo.New()
	e.GET("/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", docServerPort)))
}
