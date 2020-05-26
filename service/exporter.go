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

func (this *Exporter) Do() error {
	posts, err := this.Fetcher.Fetch()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"location": "Exporter",
		}).Error("Error during fetching data from fetcher")
		return err
	}

	var validPosts []model.Post
	for _, post := range posts {
		timePublished := post.PublishedDate.Add(time.Minute + (10 * time.Second))
		now := time.Now()

		if timePublished.Equal(now) || timePublished.After(now) {
			validPosts = append(validPosts, post)
		}
	}

	if len(validPosts) == 0 {
		logrus.Info("Nothing to do")
		return nil
	}

	logrus.Info("Found " + strconv.Itoa(len(validPosts)) + " post(s) to publish")

	err = this.Publisher.Publish(posts)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"location": "Exporter",
		}).Error("Error during publishing data")
		return err
	}
	return nil
}
