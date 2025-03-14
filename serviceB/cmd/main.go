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
	getTemp := NewGetTempoHandler()

	go func() {
        httpServer := webserver.NewWebServer(NewConfig())
		httpServer.AddHandler("GET", "/swagger/*", httpSwagger.WrapHandler)
		httpServer.AddHandler("POST", "/consulta-tempo", getTemp.HandleLabsTwo)
		fmt.Println("HTTP server running at port 8080")
		httpServer.Start()
    }()

    select {} 
}
