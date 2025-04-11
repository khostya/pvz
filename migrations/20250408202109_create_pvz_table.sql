-- +goose Up
-- +goose StatementBegin
create table if not exists pvzs
(
    id         uuid primary key not null,
    city city not null,
    registration_date timestamptz not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists pvzs;
-- +goose StatementEnd
