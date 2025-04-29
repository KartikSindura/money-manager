-- +goose Up
-- +goose StatementBegin
CREATE INDEX index_expenses_on_date ON expenses(date DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX index_expenses_on_date;
-- +goose StatementEnd
