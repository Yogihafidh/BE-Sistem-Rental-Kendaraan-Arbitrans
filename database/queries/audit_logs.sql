-- name: ListAuditLogs :many
SELECT * FROM audit_logs
ORDER BY created_at DESC;

-- name: GetAuditLogByID :one
SELECT * FROM audit_logs
WHERE id = $1
LIMIT 1;
