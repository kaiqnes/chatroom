BEGIN;

-- Drop indexes for the rooms table
DROP INDEX IF EXISTS idx_rooms_id;

-- Drop the rooms table
DROP TABLE IF EXISTS public.rooms;

COMMIT;
