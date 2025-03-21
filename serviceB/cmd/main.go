package main

import (
	"fmt"
	"labs-two-serviceb/internal/infra/web/webserver"

	_ "labs-two-serviceb/docs"

	httpSwagger "github.com/swaggo/http-swagger"
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
	httpServer.AddHandler("POST", "/consulta-tempo", getTemp.HandleLabsOne)
	fmt.Println("HTTP server running")
	httpServer.Start()
}
