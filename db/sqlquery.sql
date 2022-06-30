
create table users(
    id serial primary key,
    username varchar(50) unique not null,
    email varchar(50) unique not null,
    password varchar(50),
    age integer,
    created_at date,
    updated_at date
);

create table photo(
    id serial primary key,
    title varchar(50),
    caption varchar(50),
    photo_url varchar(50),
    user_id integer references users(id),
    ceated_date date,
    updated_at date
);

create table comment(
    id serial primary key,
    user_id integer references users(id),
    photo_id integer references photo(id),
    message varchar(50),
    created_at date,
    updated_at date
);

create table socialmedia(
    id serial primary key,
    name varchar,
    social_media_url text,
    userid integer references users(id)
);