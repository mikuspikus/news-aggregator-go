CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE table accounts (
    id SERIAL PRIMARY KEY,
    user_uid UUID not null,
    timestamp timestamp with time zone not null,
    input json,
    output json
);

create view accounts_view as
    select * from accounts order by timestamp desc;

create table news (
    id SERIAL PRIMARY KEY,
    user_uid UUID not null,
    timestamp timestamp with time zone not null,
    input json,
    output json
);
create view news_view as
    select * from news order by timestamp desc;

create table comments (
    id SERIAL PRIMARY KEY,
    user_uid UUID not null,
    timestamp timestamp with time zone not null,
    input json,
    output json
);
create view comments_view as
    select * from comments order by timestamp desc;