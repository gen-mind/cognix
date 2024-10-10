-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents ADD COLUMN IF NOT EXISTS original_url text not null default '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN IF EXISTS original_url;
-- +goose StatementEnd
