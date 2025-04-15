-- +goose Up
-- +goose StatementBegin
create index if not exists products_reception_id_idx on products using btree (reception_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists products_reception_id_idx;
-- +goose StatementEnd
