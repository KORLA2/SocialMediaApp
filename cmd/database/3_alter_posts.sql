ALTER TABLE posts add constraint posts_user_id_fkey 
FOREIGN KEY (user_id) references users(id) ON delete cascade

Alter table users add column is_active bool default false
