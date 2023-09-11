BEGIN;

-- Create the users table
CREATE TABLE IF NOT EXISTS public.users
(
    id         uuid PRIMARY KEY,
    username   VARCHAR(50)   NOT NULL UNIQUE,
    password   VARCHAR(50)   NOT NULL,
    email      VARCHAR(100)   NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for the users table
CREATE INDEX IF NOT EXISTS idx_users_id ON public.users (id);
CREATE INDEX IF NOT EXISTS idx_users_username ON public.users (username);
CREATE INDEX IF NOT EXISTS idx_users_email ON public.users (email);

COMMIT;
