BEGIN;

-- Create the messages table
CREATE TABLE IF NOT EXISTS public.messages
(
    id           uuid PRIMARY KEY,
    user_room_id uuid      NOT NULL,
    body         TEXT      NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for the messages table
CREATE INDEX IF NOT EXISTS idx_messages_id ON public.messages (id);
CREATE INDEX IF NOT EXISTS idx_user_room_id ON public.messages (user_room_id);

COMMIT;