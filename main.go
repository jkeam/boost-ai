package main

import (
	"os"

	"fmt"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	log.Debug("Start of boost ai conversation")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	baseURL := os.Getenv("BOOST_BASE_URL")
	log.Debug(fmt.Sprintf("Base url: %s", baseURL))
}
