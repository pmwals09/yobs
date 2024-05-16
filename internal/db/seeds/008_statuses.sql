-- +goose Up
-- +goose StatementBegin
INSERT INTO statuses (
    name, note,
    date,
    opportunity_id
) VALUES (
    "Applied", "",
    datetime("now", "localtime"),
    1
);

INSERT INTO statuses (
    name, note,
    date,
    opportunity_id
) VALUES (
    "Applied", "",
    datetime("now", "localtime"),
    2
);

INSERT INTO statuses (
    name, note,
    date,
    opportunity_id
) VALUES (
    "Applied", "",
    "1996-09-09 00:00:00+00:00",
    3
);

INSERT INTO statuses (
    name, note,
    date,
    opportunity_id
) VALUES (
    "Applied", "",
    "2025-11-20 00:00:00+00:00",
    4
);

INSERT INTO statuses (
    name, note,
    date,
    opportunity_id
) VALUES (
    "Applied", "",
    "2030-11-20 00:00:00+00:00",
    5
);

INSERT INTO statuses (
    name, note,
    date,
    opportunity_id
) VALUES (
    "Applied", "",
    "2023-01-01 00:00:00+00:00",
    6
);
-- +goose StatementEnd

-- +goose Down
DELETE FROM statuses;
