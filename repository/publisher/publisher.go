package publisher

import "exporter/model"

type Publisher interface {
	Publish(posts []model.Post) error
}
