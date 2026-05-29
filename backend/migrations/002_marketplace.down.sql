-- =============================================================================
-- hada-commerce: Rollback marketplace & plugin registry migration
-- Drop in reverse dependency order.
-- =============================================================================

DROP TABLE IF EXISTS plugin_installations;
DROP TABLE IF EXISTS plugin_versions;
DROP TABLE IF EXISTS plugins;
