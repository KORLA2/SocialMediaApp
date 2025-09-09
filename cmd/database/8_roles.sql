create table if not exists roles  (

id bigserial primary KEY,
name text not null unique,
level int not null default 0,
description text
);

Insert into roles (name,level,description) values('user',1,'Users can manage their own profile');
Insert into roles (name,level,description) values('moderator',2,'Moderators can Update others profiles');
Insert into roles (name,level,description) values('admin',3,'Admins can manage others profiles');