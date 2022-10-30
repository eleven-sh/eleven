package entities

import (
	"time"

	"github.com/google/uuid"
)

type Config struct {
	ID                 string              `json:"id"`
	Clusters           map[string]*Cluster `json:"clusters"`
	CreatedAtTimestamp int64               `json:"created_at_timestamp"`
}

func NewConfig() *Config {
	return &Config{
		ID:                 uuid.NewString(),
		Clusters:           map[string]*Cluster{},
		CreatedAtTimestamp: time.Now().Unix(),
	}
}
