-- +goose Up
INSERT INTO opportunity_contacts (
  opportunity_id,
  contact_id
) VALUES (
  1,
  1
);
INSERT INTO opportunity_contacts (
  opportunity_id,
  contact_id
) VALUES (
  1,
  2
);
INSERT INTO opportunity_contacts (
  opportunity_id,
  contact_id
) VALUES (
  1,
  3
);
-- +goose Down
DELETE FROM opportunity_contacts;
