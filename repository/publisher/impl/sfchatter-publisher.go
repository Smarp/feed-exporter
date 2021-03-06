package salesforce_chatter

import (
	"bytes"
	"encoding/json"
	"errors"
	"exporter/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	subjectId               = "0F94I0000009tm8SAA"
	sfApi                   = "https://smarp.my.salesforce.com/services/data/v48.0"
	chatterFeedElementsPATH = "/chatter/feed-elements"
	smarpPostUrl            = "https://salesforce.smarpshare.com/#/preview/"
)

type SFChatterPublisher struct {
}

type FeedItem struct {
	FeedElementType string       `json:"feedElementType"`
	SubjectId       string       `json:"subjectId"`
	Body            FeedItemBody `json:"body"`
}

type FeedItemBody struct {
	MessageSegments []MessageSegment `json:"messageSegments"`
}
type MessageSegment struct {
	TheType string `json:"type"`
	Text    string `json:"text"`
}

var feedItemGenerator = func(subjectId, message string) FeedItem {
	return FeedItem{
		FeedElementType: "FeedItem",
		SubjectId:       subjectId,
		Body: FeedItemBody{
			MessageSegments: []MessageSegment{{
				TheType: "Text",
				Text:    message,
			}},
		},
	}
}
var getToken = func() string {
	clientid := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	username := os.Getenv("CLIENT_USERNAME")
	pass := os.Getenv("CLIENT_PASSWORD")
	if clientid == "" || clientSecret == "" || username == "" || pass == "" {
		logrus.WithFields(logrus.Fields{}).Error("not enough credentials. can not generate token")
		return ""
	}
	url := "https://login.salesforce.com/services/oauth2/token?grant_type=password&client_id=" + clientid +
		"&client_secret=" + clientSecret + "&username=" + username + "&password=" + pass
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": url,
			"err": err,
		}).Error("Cannot Publish")
		return ""
	}

	client := &http.Client{}

	resp, err := client.Do(request)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		byteBody, _ := ioutil.ReadAll(resp.Body)
		logrus.WithFields(logrus.Fields{
			"body":        string(byteBody),
			"status code": resp.StatusCode,
		}).Error("unable to get token")
		return ""
	}
	logrus.Info("new token created")
	byteBody, _ := ioutil.ReadAll(resp.Body)

	type authRes struct {
		AccessToken string `json:"access_token"`
	}
	respBody := &authRes{}
	err = json.Unmarshal(byteBody, respBody)
	return respBody.AccessToken
}

func (SFChatterPublisher) Publish(posts []model.Post) error {
	token := getToken()
	if token == "" {
		return errors.New("no token provided")
	}
	for _, post := range posts {
		if post.Title == "" || post.Id == "" {
			logrus.WithFields(logrus.Fields{
				"title": post.Title,
				"id":    post.Id,
			}).Error("No title or id provided. Skipping")
			continue
		}
		text := post.Title + " " + smarpPostUrl + post.Id
		feedItem := feedItemGenerator(subjectId, text)

		requestByte := new(bytes.Buffer)
		json.NewEncoder(requestByte).Encode(feedItem)

		url := sfApi + chatterFeedElementsPATH
		request, err := http.NewRequest("POST", url, requestByte)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"url": url,
				"err": err,
			}).Error("Cannot Publish")
			return err
		}
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("Connection", "Keep-Alive")
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		resp, err := client.Do(request)
		defer resp.Body.Close()
		if resp.StatusCode != 201 {
			byteBody, _ := ioutil.ReadAll(resp.Body)
			logrus.WithFields(logrus.Fields{
				"body":        string(byteBody),
				"status code": resp.StatusCode,
				"title":       post.Title,
				"id":          post.Id,
			}).Error("error during publication")
		}
		logrus.WithFields(logrus.Fields{
			"id": post.Id,
		}).Info("published post")
	}
	return nil
}
