package models

import "time"

type FeatureFlag struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Key               string    `json:"key"`
	Description       string    `json:"description"`
	Enabled           bool      `json:"enabled"`
	RolloutPercentage int       `json:"rollout_percentage"`
	Environment       string    `json:"environment"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateFeatureFlagRequest struct {
	Name              string `json:"name"`
	Key               string `json:"key"`
	Description       string `json:"description"`
	Enabled           bool   `json:"enabled"`
	RolloutPercentage int    `json:"rollout_percentage"`
	Environment       string `json:"environment"`
}

type UpdateFeatureFlagRequest struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Enabled           bool   `json:"enabled"`
	RolloutPercentage int    `json:"rollout_percentage"`
	Environment       string `json:"environment"`
}
