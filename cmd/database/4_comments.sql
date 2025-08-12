CREATE TABLE  IF NOT EXISTS comments(
    id big serial primary key,
    content text NOT NULL,
    user_id bigint NOT NULL,
    post_id bigint NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(), 
)