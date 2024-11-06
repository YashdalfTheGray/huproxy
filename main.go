package main

import (
	"net/http"

	"github.com/YashdalfTheGray/huproxy/config"
	"github.com/YashdalfTheGray/huproxy/handlers"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found")
	}

	cfg, err := config.LoadConfig(log)
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	handler := handlers.NewHandler(cfg, log)

	http.HandleFunc("/ping", handler.PingHandler)
	http.HandleFunc("/page", handler.PageHandler)

	log.Info("Starting server on :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
