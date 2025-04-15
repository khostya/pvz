-- +goose Up
-- +goose StatementBegin
create index if not exists receptions_pvz_id_idx on receptions using btree (pvz_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists receptions_pvz_id_idx;
-- +goose StatementEnd
