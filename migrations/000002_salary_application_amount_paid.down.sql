ALTER TABLE salary_applications
    DROP COLUMN IF EXISTS amount,
    DROP COLUMN IF EXISTS paid_at;
