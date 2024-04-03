create table if not exists models
(
    id          text primary key,
    domain      text                                                not null,
    name        text                                                not null,
    base        text default 'runwayml/stable-diffusion-inpainting' not null,
    ckpt        bytea                                               not null,
    created_at  timestamp default now()
);

create table if not exists inferences
(
    id         text primary key,
    prompt     text,
    image      bytea not null,
    image_url  text  not null,
    mask       bytea not null,
    mask_url   text  not null,
    output     bytea not null,
    output_url text  not null,
    from_model text  not null,
    foreign key (from_model) references models (id)
);

create table if not exists base_assets
(
    id        text primary key,
    image     bytea not null,
    image_url text  not null,
    mask      bytea not null,
    mask_url  text  not null,
    domain    text  not null
);


CREATE TYPE task_variant AS ENUM ('inference', 'sample', 'finetune');

create table if not exists tasks
(
    id              text,
    source_model_id text not null,
    output_model_id text,
    task_type       task_variant not null,
    created_at      timestamp default now(),
    handled_at      timestamp,
    finished_at     timestamp,
    human_prefs     jsonb,
    prompt_embeds   bytea,
    latents         bytea,
    timesteps       bytea,
    next_latents    bytea,
    image_torchs    bytea,
    primary key (id, task_type),
    foreign key (source_model_id) references models (id),
    foreign key (output_model_id) references models (id)
);

create table if not exists assets
(
    task_id   text,
    "order"   smallint,
    prompt    text,
    image     bytea     not null,
    image_url text      not null,
    mask      bytea     not null,
    mask_url  text      not null,
    primary key (task_id, "order"),
    foreign key (task_id) references tasks (id)
);
