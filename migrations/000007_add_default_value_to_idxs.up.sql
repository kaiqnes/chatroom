BEGIN;

ALTER TABLE public.users_rooms ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE public.users ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE public.rooms ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE public.messages ALTER COLUMN id SET DEFAULT gen_random_uuid();

COMMIT;
