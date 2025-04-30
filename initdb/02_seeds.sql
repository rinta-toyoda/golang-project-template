-- Dummy data with simplified, sequential UUID-like IDs

-- 1) Users
INSERT INTO users (id, email, phone, password_hash) VALUES
  ('00000000-0000-0000-0000-000000000000', 'owner@example.com', '+14151234567', '$2a$14$DXJ0OAdIRQ6ZSme.9VJ1xeJR4vhqwfgaCatCDSWWkRn1/CHcs6H0K'),
  ('00000000-0000-0000-0000-000000000001', 'example@example.com',   '+442071838750', '$2a$14$DXJ0OAdIRQ6ZSme.9VJ1xeJR4vhqwfgaCatCDSWWkRn1/CHcs6H0K');