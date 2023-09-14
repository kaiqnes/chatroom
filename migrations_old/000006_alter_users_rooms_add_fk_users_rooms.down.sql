BEGIN;

ALTER TABLE IF EXISTS public.users_rooms DROP CONSTRAINT users_rooms_room_id_fkey;

ALTER TABLE IF EXISTS public.users_rooms DROP CONSTRAINT users_rooms_user_id_fkey;

COMMIT;