package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type services struct {
	auth string
	log string
	mail string
}
type Config struct{
	services services
}

func main() {
	app := Config{
		services: services{
			auth: "http://authentication-service",
			log: "http://logger-service",
			mail: "http://mail-service",
		},
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}