BEGIN;

TRUNCATE TABLE messages CASCADE;

TRUNCATE TABLE users_rooms CASCADE;

TRUNCATE TABLE users CASCADE;

TRUNCATE TABLE rooms CASCADE;

COMMIT;
