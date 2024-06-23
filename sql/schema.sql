create table if not exists models (
    id              text primary key,
    domain          text not null,
    name            text not null,
    parent_id       text,
    status          text not null,
    sample_id       text,
    train_id        text,
    updated_at      timestamp with time zone,
    created_at      timestamp with time zone default now()
);

insert into models (
    id, domain, name, parent_id, status, sample_id, train_id, updated_at
) values
      ('base-1', 'sessile', 'base', NULL, 'trained', '', '', now()),
      ('base-2', 'pedunculated', 'base', NULL, 'trained', '', '', now()),
      ('base-3', 'outpaint', 'base', NULL, 'trained', '', '', now());

create table if not exists samples (
    id                  text primary key,
    model_id            text not null,
    finished_at         timestamp with time zone,
    created_at          timestamp with time zone not null default now(),
    foreign key(model_id) references models(id)
);

create table if not exists trains (
    id                  text primary key,
    sample_id           text not null,
    model_id            text not null,
    created_at          timestamp with time zone not null default now(),
    finished_at         timestamp with time zone,
    foreign key(sample_id) references samples(id)
);

create table if not exists inferences (
    id              text primary key,
    model_id        text not null,
    prompt          text not null,
    neg_prompt      text not null,
    finished_at     timestamp with time zone,
    foreign key(model_id) references models(id)
);


