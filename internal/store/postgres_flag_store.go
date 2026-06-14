package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/sakshipatel29/launchguard/internal/models"
)

type PostgresFeatureFlagStore struct {
	db *sql.DB
}

func NewPostgresFeatureFlagStore(db *sql.DB) *PostgresFeatureFlagStore {
	return &PostgresFeatureFlagStore{
		db: db,
	}
}

func (s *PostgresFeatureFlagStore) Create(req models.CreateFeatureFlagRequest) (models.FeatureFlag, error) {
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

	query := `
		INSERT INTO feature_flags (
			id, name, flag_key, description, enabled, rollout_percentage, environment, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id::text, name, flag_key, description, enabled, rollout_percentage, environment, created_at, updated_at;
	`

	err := s.db.QueryRow(
		query,
		flag.ID,
		flag.Name,
		flag.Key,
		flag.Description,
		flag.Enabled,
		flag.RolloutPercentage,
		flag.Environment,
		flag.CreatedAt,
		flag.UpdatedAt,
	).Scan(
		&flag.ID,
		&flag.Name,
		&flag.Key,
		&flag.Description,
		&flag.Enabled,
		&flag.RolloutPercentage,
		&flag.Environment,
		&flag.CreatedAt,
		&flag.UpdatedAt,
	)

	if err != nil {
		if isUniqueViolation(err) {
			return models.FeatureFlag{}, ErrDuplicateFlagKey
		}
		return models.FeatureFlag{}, err
	}

	return flag, nil
}

func (s *PostgresFeatureFlagStore) List() ([]models.FeatureFlag, error) {
	query := `
		SELECT id::text, name, flag_key, description, enabled, rollout_percentage, environment, created_at, updated_at
		FROM feature_flags
		ORDER BY created_at DESC;
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	flags := []models.FeatureFlag{}

	for rows.Next() {
		var flag models.FeatureFlag

		if err := rows.Scan(
			&flag.ID,
			&flag.Name,
			&flag.Key,
			&flag.Description,
			&flag.Enabled,
			&flag.RolloutPercentage,
			&flag.Environment,
			&flag.CreatedAt,
			&flag.UpdatedAt,
		); err != nil {
			return nil, err
		}

		flags = append(flags, flag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return flags, nil
}

func (s *PostgresFeatureFlagStore) GetByID(id string) (models.FeatureFlag, error) {
	query := `
		SELECT id::text, name, flag_key, description, enabled, rollout_percentage, environment, created_at, updated_at
		FROM feature_flags
		WHERE id::text = $1;
	`

	return s.getOne(query, id)
}

func (s *PostgresFeatureFlagStore) GetByKeyAndEnvironment(key string, environment string) (models.FeatureFlag, error) {
	query := `
		SELECT id::text, name, flag_key, description, enabled, rollout_percentage, environment, created_at, updated_at
		FROM feature_flags
		WHERE flag_key = $1 AND environment = $2;
	`

	return s.getOne(query, key, environment)
}

func (s *PostgresFeatureFlagStore) Update(id string, req models.UpdateFeatureFlagRequest) (models.FeatureFlag, error) {
	query := `
		UPDATE feature_flags
		SET
			name = $1,
			description = $2,
			enabled = $3,
			rollout_percentage = $4,
			environment = $5,
			updated_at = $6
		WHERE id::text = $7
		RETURNING id::text, name, flag_key, description, enabled, rollout_percentage, environment, created_at, updated_at;
	`

	var flag models.FeatureFlag

	err := s.db.QueryRow(
		query,
		req.Name,
		req.Description,
		req.Enabled,
		req.RolloutPercentage,
		req.Environment,
		time.Now().UTC(),
		id,
	).Scan(
		&flag.ID,
		&flag.Name,
		&flag.Key,
		&flag.Description,
		&flag.Enabled,
		&flag.RolloutPercentage,
		&flag.Environment,
		&flag.CreatedAt,
		&flag.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.FeatureFlag{}, ErrFlagNotFound
		}

		if isUniqueViolation(err) {
			return models.FeatureFlag{}, ErrDuplicateFlagKey
		}

		return models.FeatureFlag{}, err
	}

	return flag, nil
}

func (s *PostgresFeatureFlagStore) Delete(id string) error {
	query := `
		DELETE FROM feature_flags
		WHERE id::text = $1;
	`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrFlagNotFound
	}

	return nil
}

func (s *PostgresFeatureFlagStore) getOne(query string, args ...interface{}) (models.FeatureFlag, error) {
	var flag models.FeatureFlag

	err := s.db.QueryRow(query, args...).Scan(
		&flag.ID,
		&flag.Name,
		&flag.Key,
		&flag.Description,
		&flag.Enabled,
		&flag.RolloutPercentage,
		&flag.Environment,
		&flag.CreatedAt,
		&flag.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.FeatureFlag{}, ErrFlagNotFound
		}

		return models.FeatureFlag{}, err
	}

	return flag, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
