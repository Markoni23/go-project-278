-- +goose Up
-- +goose StatementBegin
CREATE TABLE links (
    id BIGSERIAL PRIMARY KEY,
    original_url TEXT,
    short_name VARCHAR(255),
    short_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE links;
-- +goose StatementEnd
