package service

import (
	"exporter/model"
	"exporter/repository/fetcher"
	"exporter/repository/publisher"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Exporter struct {
	Fetcher   fetcher.Fetcher
	Publisher publisher.Publisher
}

var latestFetchTime = time.Now()

func (this *Exporter) Do() error {
	logrus.Info("checking posts from ", latestFetchTime)
	fetchedPosts, err := this.Fetcher.Fetch()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"location": "Exporter",
		}).Error("Error during fetching data from fetcher")
		return err
	}
	newLatestFetchTime := latestFetchTime
	var newPosts []model.Post
	for _, post := range fetchedPosts {
		//fmt.Println("---")
		//fmt.Println("post", post.PublishedDate.Unix())
		//fmt.Println("post", post.PublishedDate)
		//fmt.Println("latest", post.Title)
		//
		//fmt.Println("latest", latestFetchTime.Unix())
		//fmt.Println("latest", latestFetchTime)
		if post.PublishedDate.Unix() > latestFetchTime.Unix() {
			newPosts = append(newPosts, post)
			if post.PublishedDate.Unix() > newLatestFetchTime.Unix() {
				newLatestFetchTime = *post.PublishedDate
			}
		}
	}
	latestFetchTime = newLatestFetchTime
	if len(newPosts) == 0 {
		logrus.Info("Nothing to do")
		return nil
	}

	logrus.Info("Found " + strconv.Itoa(len(newPosts)) + " post(s) to publish")

	err = this.Publisher.Publish(newPosts)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"location": "Exporter",
		}).Error("Error during publishing data")
		return err
	}
	return nil
}
