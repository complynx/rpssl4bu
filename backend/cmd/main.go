package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/complynx/rpssl4bu/backend/pkg"
	gameapi "github.com/complynx/rpssl4bu/backend/pkg/GameAPI"
	"github.com/complynx/rpssl4bu/backend/pkg/game"
	"github.com/complynx/rpssl4bu/backend/pkg/random"
	"github.com/complynx/rpssl4bu/backend/pkg/server"
	"github.com/complynx/rpssl4bu/backend/pkg/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogger(logLevel, logType *string) *zap.Logger {
	// Set up logging
	var logCfg zap.Config
	if logType == nil || *logType == "text" {
		logCfg = zap.NewDevelopmentConfig()
	} else {
		logCfg = zap.NewProductionConfig()
	}
	var level zapcore.Level
	if logLevel != nil {
		var err error
		level, err = zapcore.ParseLevel(*logLevel)
		if err != nil {
			fmt.Printf("Failed to parse log level: %v\n", err)
			os.Exit(1)
		}
	} else {
		level = zap.InfoLevel
	}
	logCfg.Level.SetLevel(level)
	logger, err := logCfg.Build()
	if err != nil {
		fmt.Printf("Failed to build logger: %v\n", err)
		os.Exit(1)
	}

	return logger
}

var defaultAddr = ":8080"

func main() {
	// Parse command line arguments
	rngAddr := flag.String("rng", "", "address of the random number provider")
	addr := flag.String("addr", defaultAddr, "address and port of the server")
	logLevel := flag.String("log-level", "info", "log level (debug, info, warn, error, dpanic, panic, fatal)")
	logType := flag.String("log-type", "text", "log output type (text or json)")
	flag.Parse()

	logger := getLogger(logLevel, logType)
	defer logger.Sync()

	// Create game
	var rng pkg.RandomProvider
	if rngAddr == nil || *rngAddr == "" {
		rng = random.NewSimpleRandom("")
	} else {
		rng = random.NewProvider(*rngAddr, logger.Named("Random Provider"))
	}
	gameEngine := game.NewGame(rng)

	storage := storage.NewSimple(10)

	// Create API
	api := gameapi.NewGameAPI(gameEngine, nil, storage, logger.Named("GameAPI"))

	if addr == nil {
		addr = &defaultAddr
	}
	srv := server.StartHTTPServer(*addr, api, logger.Named("server"))

	// Wait for SIGINT or SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Gracefully shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
