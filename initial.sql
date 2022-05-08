create table user
(
    id       int primary key auto_increment,
    name     text,
    username text,
    password text,
    email    text,
    active   tinyint(1) default 1
);

create table role
(
    id   int primary key auto_increment,
    name text
);

create table user_role
(
    id      int primary key auto_increment,
    id_user int,
    id_role int,
    foreign key (id_user) references user (id),
    foreign key (id_role) references role (id)
);

create table token
(
    id            int primary key auto_increment,
    id_user       int,
    expiration    timestamp,
    refresh_token text,
    foreign key (id_user) references user (id)
);

insert into role (name)
values ('USER');
insert into role (name)
values ('ADMIN');
