CREATE TABLE posts if not exists (
    id bigserial primary key,
    title text Not null,
    content text not null,
    user_id  bigint NOT Null,
    created_at timestamp(0) with time zone Not Null Default now()

)
