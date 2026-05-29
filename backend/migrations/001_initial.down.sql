-- =============================================================================
-- hada-commerce: Rollback for 001_initial.up.sql
-- Tables are dropped in reverse dependency order (children before parents).
-- =============================================================================

DROP TABLE IF EXISTS media;
DROP TABLE IF EXISTS promos;
DROP TABLE IF EXISTS page_versions;
DROP TABLE IF EXISTS pages;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS customer_addresses;
DROP TABLE IF EXISTS customers;
