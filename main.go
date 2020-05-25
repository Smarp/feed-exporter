package main

import (
	"exporter/repository/fetcher/impl"
	salesforce_chatter "exporter/repository/publisher/impl"
	"exporter/service"
)

func main() {
	exporter := service.Exporter{
		Fetcher:   smarp.SmarpFetcher{},
		Publisher: salesforce_chatter.SFChatterPublisher{},
	}

	exporter.Do()
}
