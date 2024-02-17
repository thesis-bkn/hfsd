create table if not exists users (
    id text primary key,
    email text not null,
    password text not null,
    activated boolean default false,
    created_at bigint default date_part('epoch'::text, now()) not null
);
