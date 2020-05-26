package smarp

import (
	"encoding/json"
	"errors"
	"exporter/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

type SmarpFetcher struct {
}

func (SmarpFetcher) Fetch() ([]model.Post, error) {
	instance := os.Getenv("INSTANCE_SUBDOMAIN")
	if instance == "" {
		return nil, errors.New("no instance provided")
	}
	channelId := os.Getenv("CHANNEL_ID")
	if channelId == "" {
		return nil, errors.New("no channelId provided")
	}
	url := "https://" + instance + ".smarpshare.com/api/post3?type=published&page=0&pageSize=10&channelId=" + channelId
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"location": "Smarp Fetcher",
		}).Error("Error during fetch posts")

		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"location": "Smarp Fetcher",
		}).Error("Error during reading data")

		return nil, err
	}

	var posts []model.Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"location": "Smarp Fetcher",
		}).Error("Error during unmarshalling data")

		return nil, err
	}

	return posts, nil
}
