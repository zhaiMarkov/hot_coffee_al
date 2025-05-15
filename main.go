package main

import (
	"fmt"
	"log"
	"net/http"

	"hot-coffee/pkg/logger"

	"hot-coffee/internal/config"
	jsondb "hot-coffee/internal/dal/jsonDB"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service/usecase"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}

	logg := logger.NewLogger().GetLoggerObject(config.InfoLogPath, config.ErrorLogPath, config.DebugLogPath)
	logg.InfoLogger.Println("Configuration initialized successfully")

	repo := jsondb.NewJsonDB()
	logg.InfoLogger.Println("Initialized JSON DB repository")
	service := usecase.NewApplication(repo)
	logg.InfoLogger.Println("Application service initialized")
	handlerHTTP := handler.NewCustomHandler(service)
	logg.InfoLogger.Println("HTTP Handler created")
	handlerHTTP.LoggerINFO = logg.InfoLogger
	handlerHTTP.LoggerERROR = logg.ErrorLogger
	handlerHTTP.LoggerDEBUG = logg.DebugLogger

	router := handlerHTTP.Routing()

	log.Println(fmt.Sprintf("The server is running on port http://localhost:%d", config.Port))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
