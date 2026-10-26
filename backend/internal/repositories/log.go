package repositories

import (
	"context"
	"errors"

	"go.uber.org/zap"

	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/pkg/logger"
)

// logRepoErr logs unexpected repository failures with trace_id. Domain errors stay quiet.
func logRepoErr(ctx context.Context, op string, err error) error {
	if err == nil || isDomainErr(err) {
		return err
	}
	logger.FromContext(ctx).Error("repository error",
		zap.String("layer", "repository"),
		zap.String("op", op),
		zap.Error(err),
	)
	return err
}

func isDomainErr(err error) bool {
	return errors.Is(err, apperrors.ErrValidation) ||
		errors.Is(err, apperrors.ErrNotFound) ||
		errors.Is(err, apperrors.ErrConflict)
}
