-- +goose Up
-- +goose StatementBegin
INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    user_id
) VALUES (
    "LI",
    "Influencer",
    "Influence the people",
    "http://www.li.com",
    1
);

INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    user_id
) VALUES (
    "IL",
    "Influnencer",
    "This is a descriptions",
    "http://www.il.com",
    1
);

INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    user_id
) VALUES (
    "AB",
    "AB CEO",
    "adkjhflkajh",
    "http://www.ab.com",
    1
);

INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    user_id
) VALUES (
    "NewCo",
    "Intern",
    "interns stuff",
    "http://www.newco.m",
    1
);

INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    user_id
) VALUES (
    "A thing",
    "OK",
    "dkjfhalkjhd",
    "http://www.thing.com",
    1
);

INSERT INTO opportunities (
    company_name,
    role,
    description,
    url,
    user_id
) VALUES (
    "one more",
    "one more time",
    "oh year",
    "http://www.celebrate.com",
    1
);
-- +goose StatementEnd

-- +goose Down
DELETE FROM opportunities;
