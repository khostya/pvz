-- +goose Up
-- +goose StatementBegin
create type reception_status as enum ('close', 'in_progress');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop type if exists reception_status;
-- +goose StatementEnd
