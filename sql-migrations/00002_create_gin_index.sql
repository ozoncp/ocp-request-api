-- +goose Up
CREATE INDEX text_idx ON requests USING GIN (to_tsvector('russian', text));

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP INDEX IF EXISTS text_idx;
-- +goose StatementBegin
-- +goose StatementEnd
