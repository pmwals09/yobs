-- +goose Up
CREATE TABLE IF NOT EXISTS opportunity_contacts (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  opportunity_id INTEGER NOT NULL,
  contact_id INTEGER NOT NULL,
  FOREIGN KEY (opportunity_id) REFERENCES opportunities(id),
  FOREIGN KEY (contact_id) REFERENCES contacts(id)
);

-- +goose Down
DROP TABLE opportunity_contacts;

