-- +goose Up
CREATE TABLE IF NOT EXISTS contacts (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  company TEXT,
  title TEXT,
  phone TEXT,
  email TEXT
);

-- +goose Down
DROP TABLE contacts;
