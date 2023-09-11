BEGIN;

-- Drop indexes for the users_rooms table
DROP INDEX IF EXISTS idx_users_rooms_id;

-- Drop the users_rooms table
DROP TABLE IF EXISTS public.users_rooms;

COMMIT;
