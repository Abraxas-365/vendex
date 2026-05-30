-- Revert checkout enhancements

ALTER TABLE orders DROP COLUMN IF EXISTS cart_id;
ALTER TABLE orders DROP COLUMN IF EXISTS promo_code;
ALTER TABLE orders DROP COLUMN IF EXISTS carrier;
ALTER TABLE orders DROP COLUMN IF EXISTS tracking_number;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_method;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_status;
ALTER TABLE orders DROP COLUMN IF EXISTS billing_address;
ALTER TABLE orders DROP COLUMN IF EXISTS shipping_method;
ALTER TABLE orders DROP COLUMN IF EXISTS discount_currency;
ALTER TABLE orders DROP COLUMN IF EXISTS discount_amount;
ALTER TABLE orders DROP COLUMN IF EXISTS tax_currency;
ALTER TABLE orders DROP COLUMN IF EXISTS tax_amount;
ALTER TABLE orders DROP COLUMN IF EXISTS shipping_currency;
ALTER TABLE orders DROP COLUMN IF EXISTS shipping_amount;
ALTER TABLE orders DROP COLUMN IF EXISTS subtotal_currency;
ALTER TABLE orders DROP COLUMN IF EXISTS subtotal_amount;
