create table user_invitations(

token  bytea primary key 
user_id bigint  references users(id)
expiry timestamp(0) with time zone NOT NULL,
)