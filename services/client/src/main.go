package main

import (
	"errors"
	"os"
	"strconv"
	"os/signal"
	"syscall"
	client "github.com/7574-sistemas-distribuidos/tp-nivelador/src/client"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
)

func loadConfig() (client.ClientConfig, error) {
	agencyId := os.Getenv("AGENCY_ID")
	if agencyId == "" {
		return client.ClientConfig{}, errors.New("AGENCY_ID environment variable is required")
	}

	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		return client.ClientConfig{}, errors.New("SERVER_HOST environment variable is required")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		return client.ClientConfig{}, errors.New("SERVER_PORT environment variable is required")
	}

	batch := os.Getenv("BATCH_SIZE")
	if batch == "" {
		return client.ClientConfig{}, errors.New("BATCH_SIZE environment variable is required")
	}
	batchValue, err := strconv.Atoi(batch)
	if err != nil {
		return client.ClientConfig{}, errors.New("invalid BATCH value")
	}

	inputFile := os.Getenv("INPUT_FILE")
	if inputFile == "" {
		return client.ClientConfig{}, errors.New("INPUT_FILE environment variable is required")
	}

	outputFile := os.Getenv("OUTPUT_FILE")
	if outputFile == "" {
		return client.ClientConfig{}, errors.New("OUTPUT_FILE environment variable is required")
	}

	return client.ClientConfig{
		ServerHost: serverHost,
		ServerPort: serverPort,
		AgencyId:   agencyId,
		Batch: batchValue,
		InputFile: inputFile,
		OutputFile: outputFile,
	}, nil
}

func run(sigChan  <-chan os.Signal) int {
	config, err := loadConfig()
	if err != nil {
		logger.Error("load-config", logger.Fail, "err", err)
		return 1
	}

	client, err := client.NewClient(config)
	if err != nil {
		logger.Error("client-new", logger.Fail, "err", err)
		return 1
	}

	if err := client.Run(sigChan); err != 0 {
		logger.Error("client-run", logger.Fail, "err", err)
		return err
	}
	
	return 0
}

func main() {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGTERM)

	os.Exit(run(sigChan))
}
