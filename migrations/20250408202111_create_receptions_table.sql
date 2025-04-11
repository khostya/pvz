-- +goose Up
-- +goose StatementBegin
create table if not exists receptions
(
    id         uuid primary key not null,
    pvz_id uuid not null references pvzs(id),
    status reception_status not null,

    updated_at timestamptz,
    deleted_at timestamptz,
    created_at timestamptz      NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists receptions;
-- +goose StatementEnd
