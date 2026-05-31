-- 044 down: drop approval_requests table

DROP INDEX IF EXISTS idx_approval_requests_tenant_created;
DROP INDEX IF EXISTS idx_approval_requests_tenant_status;
DROP TABLE IF EXISTS approval_requests;
