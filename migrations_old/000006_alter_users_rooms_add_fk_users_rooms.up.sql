BEGIN;

ALTER TABLE IF EXISTS public.users_rooms ADD FOREIGN KEY (room_id) REFERENCES public.rooms (id);

ALTER TABLE IF EXISTS public.users_rooms ADD FOREIGN KEY (user_id) REFERENCES public.users (id);

COMMIT;
