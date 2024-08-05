package main

import (
	"flag"
	"fmt"
	"gwordly/config"
	"gwordly/game"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	envFile := flag.String("config", "env.yaml", "/path/to/.env")
	flag.Parse()

	appVersion := "0.1b"

	logger := config.BotLogger

	handleTermination()

	configMap := config.EnvVars(*envFile)

	if configMap.LogToFile.Enabled {
		logger.ToFile(configMap.LogToFile.Directory)
	}

	logger.SetLogLevel(configMap.LogLevel)

	fmt.Println("GWordly by Fruitgum - a console version of Wordly written on Go")
	fmt.Println("https://github.com/fruitgum/GWordly")
	fmt.Println("Version: " + appVersion)
	fmt.Println("2024")
	word := game.GetWords()
	game.Gwordly(word)
}

func handleTermination() {

	logger := config.BotLogger
	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// Register for specific signals (SIGINT and SIGTERM)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// Start a goroutine to handle signals
	go func() {
		sig := <-sigs
		logger.Error("Got signal %v", sig)
		done <- true
	}()

	// Block until a signal is received
	go func() {
		<-done
		logger.Fatal("Killed")
	}()
}
