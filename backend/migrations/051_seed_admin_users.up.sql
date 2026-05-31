-- Seed admin users with OTP enabled for demo purposes
-- These users can log in via the OTP/email flow (dev mode prints code to server logs)

INSERT INTO users (id, tenant_id, email, name, status, scopes, otp_enabled, email_verified)
VALUES
  -- Vendex Tech (tnt_demo) admins
  ('usr_demo_admin',  'tnt_demo',    'admin@vendex.ai',       'Demo Admin',     'ACTIVE', '{admin,products,orders,analytics,settings}', TRUE, TRUE),
  ('usr_demo_staff',  'tnt_demo',    'staff@vendex.ai',       'Demo Staff',     'ACTIVE', '{products,orders}',                          TRUE, TRUE),

  -- Urban Threads (tnt_fashion) admins
  ('usr_fashion_admin', 'tnt_fashion', 'admin@urbanthreads.co', 'Fashion Admin',  'ACTIVE', '{admin,products,orders,analytics,settings}', TRUE, TRUE),
  ('usr_fashion_staff', 'tnt_fashion', 'staff@urbanthreads.co', 'Fashion Staff',  'ACTIVE', '{products,orders}',                          TRUE, TRUE)
ON CONFLICT (email, tenant_id) DO UPDATE SET
  otp_enabled = TRUE,
  email_verified = TRUE,
  status = 'ACTIVE';
