-- +goose Up
CREATE TABLE IF NOT EXISTS contacts (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  company_name TEXT,
  title TEXT,
  phone TEXT,
  email TEXT,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE contacts;
