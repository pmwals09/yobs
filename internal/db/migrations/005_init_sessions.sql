-- +goose Up
CREATE TABLE IF NOT EXISTS sessions (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL UNIQUE,
  init_time DATETIME,
  expiration DATETIME NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE sessions;
