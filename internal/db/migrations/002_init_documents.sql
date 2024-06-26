-- NOTE: We don't save the document URL. We use the URL when creating something
-- clickable in order to download a document, but it's time-limited by AWS so we
-- don't try to save it with the document.

-- +goose Up
CREATE TABLE IF NOT EXISTS documents (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  file_name TEXT,
  title TEXT,
  type TEXT,
  content_type TEXT,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE documents;

