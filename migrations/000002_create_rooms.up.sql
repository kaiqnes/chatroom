BEGIN;

-- Create the rooms table
CREATE TABLE IF NOT EXISTS public.rooms
(
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(100),
    description    VARCHAR(500),
    created_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for the rooms table
CREATE INDEX IF NOT EXISTS idx_rooms_id ON public.rooms (id);

COMMIT;