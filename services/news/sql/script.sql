CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE table news (
    uid uuid PRIMARY KEY,
    user_uid UUID not null,
    title text not null,
    uri text not null,
    created_at timestamp with time zone not null,
    edited_at timestamp with time zone not null
);

create view news_view as
select * from news order by created_at desc