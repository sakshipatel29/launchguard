package evaluator

import (
	"testing"

	"github.com/sakshipatel29/launchguard/internal/models"
)

func TestEvaluateDisabledFlag(t *testing.T) {
	flag := models.FeatureFlag{
		Key:               "new_checkout_ui",
		Enabled:           false,
		RolloutPercentage: 100,
	}

	result := Evaluate(flag, "user_123")

	if result.Enabled {
		t.Fatal("expected disabled flag to evaluate to false")
	}

	if result.Reason != "feature flag is disabled" {
		t.Fatalf("expected disabled reason, got %s", result.Reason)
	}
}

func TestEvaluateZeroPercentRollout(t *testing.T) {
	flag := models.FeatureFlag{
		Key:               "new_checkout_ui",
		Enabled:           true,
		RolloutPercentage: 0,
	}

	result := Evaluate(flag, "user_123")

	if result.Enabled {
		t.Fatal("expected 0 percent rollout to evaluate to false")
	}

	if result.Reason != "rollout percentage is 0" {
		t.Fatalf("expected 0 rollout reason, got %s", result.Reason)
	}
}

func TestEvaluateHundredPercentRollout(t *testing.T) {
	flag := models.FeatureFlag{
		Key:               "new_checkout_ui",
		Enabled:           true,
		RolloutPercentage: 100,
	}

	result := Evaluate(flag, "user_123")

	if !result.Enabled {
		t.Fatal("expected 100 percent rollout to evaluate to true")
	}

	if result.Reason != "rollout percentage is 100" {
		t.Fatalf("expected 100 rollout reason, got %s", result.Reason)
	}
}

func TestEvaluateIsDeterministicForSameUser(t *testing.T) {
	flag := models.FeatureFlag{
		Key:               "payment_retry_v2",
		Enabled:           true,
		RolloutPercentage: 35,
	}

	first := Evaluate(flag, "user_123")
	second := Evaluate(flag, "user_123")

	if first.Bucket != second.Bucket {
		t.Fatalf("expected same bucket for same user and flag, got %d and %d", first.Bucket, second.Bucket)
	}

	if first.Enabled != second.Enabled {
		t.Fatalf("expected same enabled result for same user and flag")
	}
}

func TestEvaluateBucketRange(t *testing.T) {
	flag := models.FeatureFlag{
		Key:               "search_ranking_v2",
		Enabled:           true,
		RolloutPercentage: 50,
	}

	result := Evaluate(flag, "user_999")

	if result.Bucket < 1 || result.Bucket > 100 {
		t.Fatalf("expected bucket to be between 1 and 100, got %d", result.Bucket)
	}
}
