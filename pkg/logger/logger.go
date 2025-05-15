package logger

import (
	"fmt"
	"log"
	"os"
)

type CustomLogger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
	FatalLogger *log.Logger
}

func NewLogger() *CustomLogger {
	return &CustomLogger{}
}

func (LoggerObject *CustomLogger) GetLoggerObject(infoFilePath, errorFilePath, debugFilePath string) *CustomLogger {
	infoFile, err := os.OpenFile(infoFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatalf("Failed to open info log file: %v", err)
	}
	errorFile, err := os.OpenFile(errorFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}
	debugFile, err := os.OpenFile(debugFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatalf("Failed to open debug log file: %v", err)
	}

	LoggerObject.InfoLogger = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerObject.ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerObject.DebugLogger = log.New(debugFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerObject.FatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

	return LoggerObject
}

func ErrorWrapper(layer, functionName, context string, err error) error {
	return fmt.Errorf("%s %w\n", fmt.Sprintf("[Layer:%s,Function:%s,Context:%s]--->", layer, functionName, context), err)
}
