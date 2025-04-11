-- +goose Up
-- +goose StatementBegin
create type user_role as enum ('employee', 'moderator');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop type if exists user_role;
-- +goose StatementEnd
