-- +goose Up
-- +goose StatementBegin
create table if not exists users
(
    id         uuid primary key not null,
    email      varchar(500)     not null unique,
    password   varchar(2000)    not null,

    role       user_role     not null,

    updated_at timestamptz,
    deleted_at timestamptz,
    created_at timestamptz      NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
