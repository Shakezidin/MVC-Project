-- Partial SQL seed: transfer modes (no foreign key dependencies)
-- For full seed data including users with bcrypt passwords, run: make seed-go

INSERT INTO transfer_modes (code, name, description, is_active) VALUES
    ('UPI', 'UPI', 'Unified Payments Interface - instant transfer up to ₹1 lakh', true),
    ('NEFT', 'NEFT', 'National Electronic Funds Transfer - batch settlement', true),
    ('RTGS', 'RTGS', 'Real Time Gross Settlement - high value instant transfer', true),
    ('IMPS', 'IMPS', 'Immediate Payment Service - 24x7 instant transfer', true)
ON CONFLICT (code) DO NOTHING;
