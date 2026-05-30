-- Checkout enhancements: add pricing breakdown, billing address, payment, tracking, promo, cart ref

ALTER TABLE orders ADD COLUMN subtotal_amount    BIGINT        NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN subtotal_currency   CHAR(3)       NOT NULL DEFAULT 'USD';
ALTER TABLE orders ADD COLUMN shipping_amount     BIGINT        NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN shipping_currency   CHAR(3)       NOT NULL DEFAULT 'USD';
ALTER TABLE orders ADD COLUMN tax_amount          BIGINT        NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN tax_currency        CHAR(3)       NOT NULL DEFAULT 'USD';
ALTER TABLE orders ADD COLUMN discount_amount     BIGINT        NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN discount_currency   CHAR(3)       NOT NULL DEFAULT 'USD';
ALTER TABLE orders ADD COLUMN shipping_method     VARCHAR(255)  NULL;
ALTER TABLE orders ADD COLUMN billing_address     JSONB         NULL;
ALTER TABLE orders ADD COLUMN payment_status      VARCHAR(50)   NOT NULL DEFAULT 'pending';
ALTER TABLE orders ADD COLUMN payment_method      VARCHAR(100)  NULL;
ALTER TABLE orders ADD COLUMN tracking_number     VARCHAR(255)  NULL;
ALTER TABLE orders ADD COLUMN carrier             VARCHAR(100)  NULL;
ALTER TABLE orders ADD COLUMN promo_code          VARCHAR(100)  NULL;
ALTER TABLE orders ADD COLUMN cart_id             VARCHAR(36)   NULL;
