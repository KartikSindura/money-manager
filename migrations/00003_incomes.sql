-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS incomes (
  id BIGSERIAL PRIMARY KEY,
  -- user_id
  amount FLOAT NOT NULL,
  category_id INT REFERENCES categories(id),
  source TEXT,
  note TEXT,
  date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incomes;
-- +goose StatementEnd
