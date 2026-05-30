-- Migration 013 rollback: Drop payment and refund tables

DROP TABLE IF EXISTS refunds;
DROP TABLE IF EXISTS payments;
