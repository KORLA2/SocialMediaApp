ALTER table comments add constraint comments_user_id_fkey 
FOREIGN KEY (user_id) references users(id) ON Delete cascade;

ALTER table comments add constraint comments_post_id_fkey 
FOREIGN KEY (post_id) references post(id) ON Delete cascade;