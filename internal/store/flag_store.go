package store

import (
	"errors"

	"github.com/sakshipatel29/launchguard/internal/models"
)

var (
	ErrFlagNotFound     = errors.New("feature flag not found")
	ErrDuplicateFlagKey = errors.New("feature flag key already exists")
)

type FeatureFlagStore interface {
	Create(req models.CreateFeatureFlagRequest) (models.FeatureFlag, error)
	List() ([]models.FeatureFlag, error)
	GetByID(id string) (models.FeatureFlag, error)
	GetByKeyAndEnvironment(key string, environment string) (models.FeatureFlag, error)
	Update(id string, req models.UpdateFeatureFlagRequest) (models.FeatureFlag, error)
	Delete(id string) error
}
