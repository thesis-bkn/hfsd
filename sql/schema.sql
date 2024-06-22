create table if not exists samples (
    id                  text primary key,
    model_id            text not null,
    finished_at         timestamp with time zone,
    created_at          timestamp with time zone not null default now(),
    foreign key(sampling_model_id) references models(id)
);

create table if not exists trains (
    id                  text primary key,
    sample_id           text not null,
    model_id            text not null,
    created_at          timestamp with time zone not null default now(),
    finished_at         timestamp with time zone,
    foreign key(sample_id) references samples(id)
);

create table if not exists models (
    id              text primary key,
    domain          text not null,
    name            text not null,
    parent_id       text,
    status          text not null,
    sample_id       text,
    train_id        text,
    updated_at      timestamp with time zone,
    created_at      timestamp with time zone default now(),
    foreign key (parent_id) references models(id),
    foreign key (sample_id) references samples(id),
    foreign key (train_id) references trains(id)
);

create table if not exists inferences (
    id              text primary key,
    model_id        text not null,
    prompt          text not null,
    neg_prompt      text not null,
    finished_at     timestamp with time zone,
    foreign key(model_id) references models(id)
);
