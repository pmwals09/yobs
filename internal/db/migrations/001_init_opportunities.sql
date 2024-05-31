-- +goose Up
CREATE TABLE IF NOT EXISTS opportunities (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  company_name TEXT,
  role TEXT,
  description TEXT,
  url TEXT,
  application_date TEXT,
  status TEXT,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE opportunities;
