package main

import (
	"fmt"
	_ "labs-two-service-a/docs"
	"labs-two-service-a/internal/infra/web/webserver"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Consulta CEP
// @version 1.0
// @description Tudo Azul Temperaturas
// @BasePath /
func main() {
	getTemp := NewGetCepHandler()

	go func() {
        httpServer := webserver.NewWebServer(NewConfig())
		httpServer.AddHandler("GET", "/swagger/*", httpSwagger.WrapHandler)
		httpServer.AddHandler("POST", "/consulta-cep", getTemp.HandleLabsTwo)
		fmt.Println("HTTP server is running")
		httpServer.Start()
    }()

    select {} 
}
