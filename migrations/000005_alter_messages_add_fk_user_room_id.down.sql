BEGIN;

ALTER TABLE IF EXISTS public.messages DROP CONSTRAINT messages_user_id_fkey;
ALTER TABLE IF EXISTS public.messages DROP CONSTRAINT messages_room_id_fkey;

COMMIT;