-- +goose Up
-- +goose StatementBegin
create index if not exists receptions_date_time_idx on receptions using btree (date_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists receptions_date_time_idx;
-- +goose StatementEnd
