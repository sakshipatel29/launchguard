package store

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/sakshipatel29/launchguard/internal/models"
)

var (
	ErrFlagNotFound     = errors.New("feature flag not found")
	ErrDuplicateFlagKey = errors.New("feature flag key already exists")
)

type FeatureFlagStore struct {
	mu    sync.RWMutex
	flags map[string]models.FeatureFlag
}

func NewFeatureFlagStore() *FeatureFlagStore {
	return &FeatureFlagStore{
		flags: make(map[string]models.FeatureFlag),
	}
}

func (s *FeatureFlagStore) Create(req models.CreateFeatureFlagRequest) (models.FeatureFlag, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, flag := range s.flags {
		if flag.Key == req.Key && flag.Environment == req.Environment {
			return models.FeatureFlag{}, ErrDuplicateFlagKey
		}
	}

	now := time.Now().UTC()

	flag := models.FeatureFlag{
		ID:                uuid.NewString(),
		Name:              req.Name,
		Key:               req.Key,
		Description:       req.Description,
		Enabled:           req.Enabled,
		RolloutPercentage: req.RolloutPercentage,
		Environment:       req.Environment,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	s.flags[flag.ID] = flag

	return flag, nil
}

func (s *FeatureFlagStore) List() []models.FeatureFlag {
	s.mu.RLock()
	defer s.mu.RUnlock()

	flags := make([]models.FeatureFlag, 0, len(s.flags))

	for _, flag := range s.flags {
		flags = append(flags, flag)
	}

	return flags
}

func (s *FeatureFlagStore) GetByID(id string) (models.FeatureFlag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	flag, exists := s.flags[id]
	if !exists {
		return models.FeatureFlag{}, ErrFlagNotFound
	}

	return flag, nil
}

func (s *FeatureFlagStore) Update(id string, req models.UpdateFeatureFlagRequest) (models.FeatureFlag, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	flag, exists := s.flags[id]
	if !exists {
		return models.FeatureFlag{}, ErrFlagNotFound
	}

	flag.Name = req.Name
	flag.Description = req.Description
	flag.Enabled = req.Enabled
	flag.RolloutPercentage = req.RolloutPercentage
	flag.Environment = req.Environment
	flag.UpdatedAt = time.Now().UTC()

	s.flags[id] = flag

	return flag, nil
}

func (s *FeatureFlagStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.flags[id]
	if !exists {
		return ErrFlagNotFound
	}

	delete(s.flags, id)

	return nil
}
