-- +goose Up
-- +goose StatementBegin
INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    application_date,
    status,
    user_id
) VALUES (
    "LI",
    "Influencer",
    "Influence the people",
    "http://www.li.com",
    datetime("now", "localtime"),
    "Applied",
    1
);

INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    application_date,
    status,
    user_id
) VALUES (
    "IL",
    "Iflunencer",
    "This is a descriptions",
    "http://www.il.com",
    datetime("now", "localtime"),
    "Applied",
    1
);

INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    application_date,
    status,
    user_id
) VALUES (
    "AB",
    "AB CEO",
    "adkjhflkajh",
    "http://www.ab.com",
    "1996-09-09 00:00:00+00:00",
    "Applied",
    1
);

INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    application_date,
    status,
    user_id
) VALUES (
    "NewCo",
    "Intern",
    "interns stuff",
    "http://www.newco.m",
    "2025-11-20 00:00:00+00:00",
    "Applied",
    1
);

INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    application_date,
    status,
    user_id
) VALUES (
    "A thing",
    "OK",
    "dkjfhalkjhd",
    "http://www.thing.com",
    "2030-11-20 00:00:00+00:00",
    "Applied",
    1
);

INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    application_date,
    status,
    user_id
) VALUES (
    "one more",
    "one more time",
    "oh year",
    "http://www.celebrate.com",
    "2023-01-01 00:00:00+00:00",
    "Applied",
    1
);
-- +goose StatementEnd

-- +goose Down
DELETE FROM opportunities;
