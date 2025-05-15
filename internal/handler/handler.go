package handler

import (
	"log"

	"hot-coffee/internal/service"
)

type CustomHandler struct {
	Service     service.ServiceModule
	LoggerINFO  *log.Logger
	LoggerERROR *log.Logger
	LoggerDEBUG *log.Logger
	LoggerWARN  *log.Logger
}

func NewCustomHandler(serviceObject service.ServiceModule) *CustomHandler {
	return &CustomHandler{Service: serviceObject}
}
