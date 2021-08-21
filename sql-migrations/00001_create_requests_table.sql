-- +goose Up
CREATE TABLE requests
(
    id      SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    type    BIGINT NOT NULL,
    text    TEXT NOT NULL
);

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS requests;
-- +goose StatementBegin
-- +goose StatementEnd
