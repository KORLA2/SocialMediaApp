create table if not exists followers(
    user_id bigint not null ,
    follower_id bigint not null,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT now(),
    primary key (user_id,follower_id), --composite-key
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (follower_id) references users(id) on delete cascade
)