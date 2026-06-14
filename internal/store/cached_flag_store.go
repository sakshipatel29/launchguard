package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/sakshipatel29/launchguard/internal/models"
)

type CachedFeatureFlagStore struct {
	primary FeatureFlagStore
	redis   *redis.Client
	ttl     time.Duration
}

func NewCachedFeatureFlagStore(primary FeatureFlagStore, redisClient *redis.Client, ttl time.Duration) *CachedFeatureFlagStore {
	return &CachedFeatureFlagStore{
		primary: primary,
		redis:   redisClient,
		ttl:     ttl,
	}
}

func (s *CachedFeatureFlagStore) Create(req models.CreateFeatureFlagRequest) (models.FeatureFlag, error) {

	flag, err := s.primary.Create(req)
	if err != nil {
		return models.FeatureFlag{}, err
	}

	s.setFlagCache(context.Background(), flag)

	return flag, nil
}

func (s *CachedFeatureFlagStore) List() ([]models.FeatureFlag, error) {
	return s.primary.List()
}

func (s *CachedFeatureFlagStore) GetByID(id string) (models.FeatureFlag, error) {
	return s.primary.GetByID(id)
}

func (s *CachedFeatureFlagStore) GetByKeyAndEnvironment(key string, environment string) (models.FeatureFlag, error) {
	ctx := context.Background()
	cacheKey := buildFlagCacheKey(key, environment)

	cachedFlag, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var flag models.FeatureFlag

		if jsonErr := json.Unmarshal([]byte(cachedFlag), &flag); jsonErr == nil {
			return flag, nil
		}
	}

	flag, err := s.primary.GetByKeyAndEnvironment(key, environment)
	if err != nil {
		return models.FeatureFlag{}, err
	}

	s.setFlagCache(ctx, flag)

	return flag, nil
}

func (s *CachedFeatureFlagStore) Update(id string, req models.UpdateFeatureFlagRequest) (models.FeatureFlag, error) {
	oldFlag, _ := s.primary.GetByID(id)

	flag, err := s.primary.Update(id, req)
	if err != nil {
		return models.FeatureFlag{}, err
	}

	ctx := context.Background()

	if oldFlag.Key != "" {
		s.redis.Del(ctx, buildFlagCacheKey(oldFlag.Key, oldFlag.Environment))
	}

	s.setFlagCache(ctx, flag)

	return flag, nil
}

func (s *CachedFeatureFlagStore) Delete(id string) error {
	oldFlag, _ := s.primary.GetByID(id)

	err := s.primary.Delete(id)
	if err != nil {
		return err
	}

	if oldFlag.Key != "" {
		s.redis.Del(context.Background(), buildFlagCacheKey(oldFlag.Key, oldFlag.Environment))
	}

	return nil
}

func (s *CachedFeatureFlagStore) setFlagCache(ctx context.Context, flag models.FeatureFlag) {
	data, err := json.Marshal(flag)
	if err != nil {
		return
	}

	cacheKey := buildFlagCacheKey(flag.Key, flag.Environment)

	if err := s.redis.Set(ctx, cacheKey, data, s.ttl).Err(); err != nil {
		return
	}
}

func buildFlagCacheKey(key string, environment string) string {
	return fmt.Sprintf("flag:%s:%s", environment, key)
}
