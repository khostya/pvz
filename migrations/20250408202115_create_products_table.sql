-- +goose Up
-- +goose StatementBegin
create table if not exists products
(
    id         uuid primary key not null,
    type product_type not null,
    reception_id uuid not null references receptions(id),
    date_time timestamptz not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists products;
-- +goose StatementEnd
