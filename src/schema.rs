// @generated automatically by Diesel CLI.

diesel::table! {
    users (id) {
        id -> Text,
        name -> Nullable<Text>,
        password -> Nullable<Text>,
        email -> Nullable<Text>,
        activated -> Nullable<Bool>,
        created_at -> Int8,
    }
}
