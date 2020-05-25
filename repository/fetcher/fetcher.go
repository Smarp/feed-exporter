package fetcher

import "exporter/model"

type Fetcher interface {
	Fetch() ([]model.Post, error)
}
