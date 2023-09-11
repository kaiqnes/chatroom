BEGIN;

-- Drop indexes for the users table
DROP INDEX IF EXISTS idx_users_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

-- Drop the users table
DROP TABLE IF EXISTS public.users;

COMMIT;
