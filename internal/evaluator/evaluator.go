package evaluator

import (
	"fmt"
	"hash/fnv"

	"github.com/sakshipatel29/launchguard/internal/models"
)

type EvaluationResult struct {
	Enabled bool
	Bucket  int
	Reason  string
}

func Evaluate(flag models.FeatureFlag, userID string) EvaluationResult {
	if !flag.Enabled {
		return EvaluationResult{
			Enabled: false,
			Bucket:  calculateBucket(flag.Key, userID),
			Reason:  "feature flag is disabled",
		}
	}

	if flag.RolloutPercentage <= 0 {
		return EvaluationResult{
			Enabled: false,
			Bucket:  calculateBucket(flag.Key, userID),
			Reason:  "rollout percentage is 0",
		}
	}

	if flag.RolloutPercentage >= 100 {
		return EvaluationResult{
			Enabled: true,
			Bucket:  calculateBucket(flag.Key, userID),
			Reason:  "rollout percentage is 100",
		}
	}

	bucket := calculateBucket(flag.Key, userID)

	if bucket <= flag.RolloutPercentage {
		return EvaluationResult{
			Enabled: true,
			Bucket:  bucket,
			Reason:  "user included in rollout",
		}
	}

	return EvaluationResult{
		Enabled: false,
		Bucket:  bucket,
		Reason:  fmt.Sprintf("user bucket %d is outside rollout percentage %d", bucket, flag.RolloutPercentage),
	}
}

func calculateBucket(flagKey string, userID string) int {
	h := fnv.New32a()
	h.Write([]byte(flagKey + ":" + userID))

	return int(h.Sum32()%100) + 1
}
