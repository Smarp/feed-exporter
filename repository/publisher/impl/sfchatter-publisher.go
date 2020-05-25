package salesforce_chatter

import "exporter/model"

type SFChatterPublisher struct {
}

func (SFChatterPublisher) Publish(posts []model.Post) error {
	return nil
}
