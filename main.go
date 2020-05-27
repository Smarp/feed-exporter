package main

import (
	"exporter/repository/fetcher/impl"
	salesforce_chatter "exporter/repository/publisher/impl"
	"exporter/service"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	exporter := service.Exporter{
		Fetcher:   smarp.SmarpFetcher{},
		Publisher: salesforce_chatter.SFChatterPublisher{},
	}

	for {
		logrus.Info("starting job...")
		exporter.Do()
		logrus.Info("finishing job.")
		time.Sleep(60 * time.Second)
	}
}
