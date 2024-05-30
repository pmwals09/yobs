-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS statuses (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  note TEXT,
  date DATETIME,
  opportunity_id INTEGER NOT NULL,
  FOREIGN KEY (opportunity_id) REFERENCES opportunities(id)
  ON DELETE CASCADE
);

INSERT INTO statuses (
  name,
  note,
  date,
  opportunity_id
) SELECT 
  status, 
  "" AS note,
  application_date,
  id
FROM opportunities;

CREATE TABLE IF NOT EXISTS opportunities_old (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  company_name TEXT,
  role TEXT,
  description TEXT,
  url TEXT,
  application_date DATETIME,
  status TEXT,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
INSERT INTO opportunities_old SELECT * FROM opportunities;
DROP TABLE opportunities;
CREATE TABLE IF NOT EXISTS opportunities (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  company_name TEXT,
  role TEXT,
  description TEXT,
  url TEXT,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
INSERT INTO opportunities (
  id,
  company_name,
  role,
  description,
  url,
  user_id
) SELECT id, company_name, role, description, url, user_id FROM opportunities_old;
DROP TABLE opportunities_old;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS opportunities_old (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  company_name TEXT,
  role TEXT,
  description TEXT,
  url TEXT,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
INSERT INTO opportunities_old (
  id,
  company_name,
  role,
  description,
  url,
  user_id
) SELECT id, company_name, role, description, url, user_id FROM opportunities;
DROP TABLE opportunities;
CREATE TABLE IF NOT EXISTS opportunities (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  company_name TEXT,
  role TEXT,
  description TEXT,
  url TEXT,
  application_date DATETIME,
  status TEXT,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
INSERT INTO opportunities (
  id,
  company_name,
  role,
  description,
  url,
  user_id
) SELECT id, company_name, role, description, url, user_id FROM opportunities_old;
DROP TABLE opportunities_old;
INSERT INTO opportunities (
  application_date,
  status
) SELECT
  application_date,
  status
FROM statuses
FULL OUTER JOIN opportunities ON statuses.opportunity_id = opportunities.id;
DROP TABLE statuses;
-- +goose StatementEnd
