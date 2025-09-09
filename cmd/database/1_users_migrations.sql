CREATE EXTENSION IF NOT EXISTS citext;


CREATE TABLE   IF NOT EXISTS users(
    id bigserial primary key,
    user_name citext NOT NULL UNIQUE,
    email citext NOT NULL UNIQUE ,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
)

