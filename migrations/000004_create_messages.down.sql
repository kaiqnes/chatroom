BEGIN;

-- Drop indexes for the messages table
DROP INDEX IF EXISTS idx_messages_id;
DROP INDEX IF EXISTS idx_user_room_id;

-- Drop the messages table
DROP TABLE IF EXISTS public.messages;

COMMIT;
