package models

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

type LogInterface interface {
	Create(error) error
}

type Log struct {
	message   string
	timeStamp time.Time
}

type DBLogger struct {
	DB *sql.DB
}

type MockLogger struct {
	ErrorLog []error
}

func (logger *DBLogger) Create(err error) error {
	errorTime := time.Now()
	_, logError := logger.DB.Exec(`
	INSERT INTO logs(message, timestamp) 
	VALUES ($1, $2)
	`, err.Error(), errorTime)
	if logError != nil {
		return logError
	}
	return nil
}

func (ml *MockLogger) Create(err error) error {
	ml.ErrorLog = append(ml.ErrorLog, err)
	return nil
}

func LoggerMiddleware(loggerInterface LogInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					loggerInterface.Create(fmt.Errorf("panic: %v", err))
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			ctx := context.WithValue(r.Context(), "logger", loggerInterface)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
