package main

import (
	"fmt"
	_ "labs-two-service-b/docs"
	"labs-two-service-b/internal/infra/web/webserver"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Consulta Tempo
// @version 1.0
// @description Consulta Temperaturas
func main() {
	tracingProvider, cleanup := InitializeTracing()
	defer cleanup()

	getTemp := NewGetTempoHandler()

	go func() {
        httpServer := webserver.NewWebServer(NewConfig(), tracingProvider)
		httpServer.AddHandler("GET", "/swagger/*", httpSwagger.WrapHandler)
		httpServer.AddHandler("POST", "/consulta-tempo", getTemp.HandleLabsTwo)
		fmt.Println("HTTP server is running")
		httpServer.Start()
    }()

    select {} 
}
