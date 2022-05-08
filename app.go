package main

import (
	"github.com/ffmoyano/goApi/database"
	"github.com/ffmoyano/goApi/logger"

	"net/http"
	"os"
	"time"
)

func main() {

	defer logger.Close()

	database.Open()
	defer database.Close()

	httpRouter := http.NewServeMux()
	addHandlers(httpRouter)

	server := &http.Server{
		Addr:         os.Getenv("port"),
		Handler:      httpRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	logger.InfoLogger.Printf("Running on port %s", os.Getenv("port"))
	if err := server.ListenAndServe(); err != nil {
		logger.ErrorLogger.Panicf("Couldn't start server:  %s", err)
	}

}
