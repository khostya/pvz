-- +goose Up
-- +goose StatementBegin
create type city as enum ('Москва', 'Санкт-Петербург', 'Казань');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop type if exists city;
-- +goose StatementEnd