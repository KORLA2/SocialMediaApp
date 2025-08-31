alter table users add colum role_id bigint references roles(id)

Update users set role_id=(
    select if from roles where name="user"
)

alter table users alter column role_id set not null