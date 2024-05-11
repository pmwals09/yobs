-- +goose Up
-- +goose StatementBegin
INSERT INTO contacts (
    name,
    company_name,
    title,
    phone,
    email,
    user_id
) VALUES (
    "Name LastName",
    "BigCorp",
    "Lorge Title",
    "123-456-7890",
    "email@bigcorp.com",
    1
);
INSERT INTO contacts (
    name,
    company_name,
    title,
    phone,
    email,
    user_id
) VALUES (
    "Peter Gibbons",
    "Initech",
    "Senior TPS Manager",
    "890-567-1234",
    "pgibbons@initech.com",
    1
);
INSERT INTO contacts (
    name,
    company_name,
    title,
    phone,
    email,
    user_id
) VALUES (
    "Dante Hicks",
    "Quick Stop Groceries",
    "Retail Clerk",
    "567-890-1234",
    "dante@quickstop.com",
    1
);
-- +goose StatementEnd

-- +goose Down
DELETE FROM contacts;
