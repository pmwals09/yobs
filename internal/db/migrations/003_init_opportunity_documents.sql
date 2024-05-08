-- +goose Up
CREATE TABLE IF NOT EXISTS opportunity_documents (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  opportunity_id INTEGER NOT NULL,
  document_id INTEGER NOT NULL,
  FOREIGN KEY (opportunity_id) REFERENCES opportunities(id),
  FOREIGN KEY (document_id) REFERENCES documents(id)
);

-- +goose Down
DROP TABLE opportunity_documents;

