-- +goose Up
-- +goose StatementBegin
create index if not exists users_email_idx on users using btree (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists users_email_idx;
-- +goose StatementEnd
