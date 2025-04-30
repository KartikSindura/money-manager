-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories
ADD CONSTRAINT unique_user_category UNIQUE (user_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories
DROP CONSTRAINT unique_user_category;
-- +goose StatementEnd
