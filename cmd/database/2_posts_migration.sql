CREATE TABLE  if not exists posts (
    id bigserial primary key,
    title text Not null,
    content text not null,
    user_id  bigint NOT Null,
    created_at timestamp(0) with time zone Not Null Default now(),
    updated_at timestamp(0) with time zone Not Null Default now(),
    tags varchar(200)[]
)
