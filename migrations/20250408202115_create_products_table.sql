-- +goose Up
-- +goose StatementBegin
create table if not exists products
(
    id         uuid primary key not null,
    type product_type not null,
    pvz_id uuid not null references pvzs(id),
    reception_id uuid not null references receptions(id),

    updated_at timestamptz,
    deleted_at timestamptz,
    created_at timestamptz      NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists products;
-- +goose StatementEnd
