-- Your SQL goes here
create table if not exists users (
    id text primary key,
    name text,
    password text,
    email text unique,
    activated boolean,
    created_at BIGINT DEFAULT DATE_PART('epoch'::TEXT, NOW()) NOT NULL
);
