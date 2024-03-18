// @generated automatically by Diesel CLI.

diesel::table! {
    assets (task_id, order) {
        task_id -> Text,
        order -> Int2,
        image -> Nullable<Bytea>,
        image_url -> Nullable<Text>,
        mask -> Nullable<Bytea>,
        mask_url -> Nullable<Text>,
    }
}

diesel::table! {
    base_assets (id) {
        id -> Text,
        image -> Nullable<Bytea>,
        image_url -> Nullable<Text>,
        mask -> Nullable<Bytea>,
        mask_url -> Nullable<Text>,
        domain -> Nullable<Text>,
    }
}

diesel::table! {
    inferences (id) {
        id -> Text,
        image -> Nullable<Bytea>,
        image_url -> Nullable<Text>,
        mask -> Nullable<Bytea>,
        mask_url -> Nullable<Text>,
        output -> Nullable<Bytea>,
        output_url -> Nullable<Text>,
        from_model -> Nullable<Text>,
    }
}

diesel::table! {
    models (id) {
        id -> Text,
        domain -> Nullable<Text>,
        name -> Nullable<Text>,
        base -> Nullable<Text>,
        ckpt -> Nullable<Bytea>,
    }
}

diesel::table! {
    tasks (id) {
        id -> Text,
        source_model_id -> Nullable<Text>,
        output_model_id -> Nullable<Text>,
        task_type -> Nullable<Text>,
        created_at -> Nullable<Timestamp>,
        handled_at -> Nullable<Timestamp>,
        finished_at -> Nullable<Timestamp>,
        human_prefs -> Nullable<Jsonb>,
        prompt_embeds -> Nullable<Bytea>,
        latents -> Nullable<Bytea>,
        timesteps -> Nullable<Bytea>,
        next_latents -> Nullable<Bytea>,
        image_torchs -> Nullable<Bytea>,
    }
}

diesel::table! {
    users (id) {
        id -> Text,
        email -> Text,
        password -> Text,
        activated -> Nullable<Bool>,
        created_at -> Int8,
    }
}

diesel::joinable!(assets -> tasks (task_id));
diesel::joinable!(inferences -> models (from_model));

diesel::allow_tables_to_appear_in_same_query!(
    assets,
    base_assets,
    inferences,
    models,
    tasks,
    users,
);
