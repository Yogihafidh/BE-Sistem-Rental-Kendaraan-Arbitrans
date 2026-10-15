package audit

import (
	"context"
	"encoding/json"

	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func Create(ctx context.Context, queries *sqlc.Queries, userID pgtype.UUID, action, entityType string, entityID pgtype.UUID, value any) error {
	// 1. Convert value to JSON
	payload, _ := json.Marshal(value)

	// 2. Create audit log
	_, err := queries.CreateAuditLog(ctx, sqlc.CreateAuditLogParams{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		NewValue:   payload,
	})

	// 3. Return error if any
	return err
}
