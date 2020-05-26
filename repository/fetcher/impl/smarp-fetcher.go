package smarp

import (
	"encoding/json"
	"exporter/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type SmarpFetcher struct {
}

func (SmarpFetcher) Fetch() ([]model.Post, error) {
	resp, err := http.Get("https://salesforce.smarpshare.com/api/post3?type=published&page=0&pageSize=10&channelId=5ecb8e663792c41ffe9f79de")
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
