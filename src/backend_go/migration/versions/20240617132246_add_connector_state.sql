-- +goose Up
-- +goose StatementBegin
ALTER TABLE connectors ADD COLUMN IF NOT EXISTS state jsonb not null default '{}'::jsonb;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE connectors DROP COLUMN IF EXISTS state;
-- +goose StatementEnd
