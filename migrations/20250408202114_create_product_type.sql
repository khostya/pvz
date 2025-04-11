-- +goose Up
-- +goose StatementBegin
create type product_type as enum ('электроника', 'одежда', 'обувь');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop type if exists product_type;
-- +goose StatementEnd