BEGIN;

ALTER TABLE IF EXISTS public.messages ADD FOREIGN KEY (user_room_id) REFERENCES public.users_rooms (id);

COMMIT;