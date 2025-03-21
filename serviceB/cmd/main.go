package main

import (
	"fmt"
	"labs-two-serviceb/internal/infra/web/webserver"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "labs-two-serviceb/docs" 
)

// @title Tudo Azul API
// @version 1.0
// @description Tudo Azul Temperaturas
// @BasePath /
func main() {
	tracingProvider, cleanup := InitializeTracing()
	defer cleanup()

	getTemp := NewGetTempoHandler()

	httpServer := webserver.NewWebServer(NewConfig(), tracingProvider)
	httpServer.AddHandler("GET", "/swagger/*", httpSwagger.WrapHandler)
	httpServer.AddHandler("GET", "/get-temp", getTemp.HandleLabsOne)
	fmt.Println("HTTP server running at port 8080")
	httpServer.Start()
}
