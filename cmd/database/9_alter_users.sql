alter table users add column role_id bigint not null references roles(id);

Update users set role_id=(
    select id from roles where name='user'
);

