CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE table comments (
    id SERIAL PRIMARY KEY,
    user_uid UUID not null,
    news_uid UUID not null,
    body text not null,
    created_at timestamp with time zone not null,
    edited_at timestamp with time zone not null
);

create index comments_index on comments using btree(news_uid);

create view comments_view as
    select * from comments order by created_at desc