BEGIN;

-- Create the users_rooms table
CREATE TABLE IF NOT EXISTS public.users_rooms
(
    id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id   uuid      NOT NULL,
    room_id   uuid      NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    left_at   TIMESTAMP
);

-- Create indexes for the users_rooms table
CREATE INDEX IF NOT EXISTS idx_users_rooms_id ON public.users_rooms (id);

COMMIT;