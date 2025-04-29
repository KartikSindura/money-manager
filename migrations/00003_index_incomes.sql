-- +goose Up
-- +goose StatementBegin
CREATE INDEX index_incomes_on_date ON incomes(date DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX index_incomes_on_date;
-- +goose StatementEnd
