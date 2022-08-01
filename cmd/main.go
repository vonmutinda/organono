package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/logger"
	"github.com/vonmutinda/organono/app/utils"
	"github.com/vonmutinda/organono/app/web/router"
)

const defaultPort = "3000"

func main() {

	logger.Initialize()
	defer logger.Flush()

	var envFilePath string

	flag.StringVar(&envFilePath, "e", "", "Path to .env file")
	flag.Parse()

	if envFilePath != "" {
		err := godotenv.Load(envFilePath)
		if err != nil {
			logger.Fatalf("Failed to load env file err = %v", err)
		}
	}

	logger.Infof("ENVIRONMENT=[%v]", os.Getenv("ENVIRONMENT"))

	dB := db.InitDB()
	defer dB.Close()

	utils.LoadTestData(dB)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router.BuildRouter(dB),
	}

	done := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		logger.Info("Process terminated...shutting down")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server shut down error = %v", err)
		}

		close(done)
	}()

	logger.Infof("Server starting... Listening on port :%v", port)

	err := server.ListenAndServe()
	if err != nil {
		switch err {
		case http.ErrServerClosed:
			logger.Info("Server successfully shut down!")
		default:
			logger.Fatal("Server shut down unexpectedly!")
		}
	}
}
