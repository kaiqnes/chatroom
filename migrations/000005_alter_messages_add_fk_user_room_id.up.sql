BEGIN;

ALTER TABLE IF EXISTS public.messages ADD FOREIGN KEY (user_id) REFERENCES public.users (id);
ALTER TABLE IF EXISTS public.messages ADD FOREIGN KEY (room_id) REFERENCES public.rooms (id);

COMMIT;