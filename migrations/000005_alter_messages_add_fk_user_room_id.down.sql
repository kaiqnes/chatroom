BEGIN;

ALTER TABLE IF EXISTS public.messages DROP CONSTRAINT messages_user_room_id_fkey;

COMMIT;