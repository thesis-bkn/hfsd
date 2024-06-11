create table if not exists samples (
    save_dir        timestamp with time zone not null primary key,
    resume_from     timestamp with time zone,
    image_fn        text not null,
    prompt_fn       text not null,
    created_at      timestamp with time zone not null default now(),
    finished_at     timestamp with time zone
);

create table if not exists ratings (
    json_path       timestamp with time zone not null primary key,
    sample_id       timestamp with time zone,
    created_at      timestamp with time zone not null default now(),
    foreign key(sample_id) references samples(save_dir)
);

create table if not exists trains (
    log_path        timestamp with time zone not null primary key,
    rating_id       timestamp with time zone,
    created_at      timestamp with time zone not null default now(),
    finished_at     timestamp with time zone,
    foreign key(rating_id) references ratings(json_path)
);

create table if not exists inferences (
    output_path     timestamp with time zone not null primary key,
    resume_from     timestamp with time zone,
    image_path      text not null,
    mask_path       text not null,
    prompt          text not null,
    neg_prompt      text not null,
    finished_at     timestamp with time zone,
    foreign key(rating_id) references ratings(json_path)
);
