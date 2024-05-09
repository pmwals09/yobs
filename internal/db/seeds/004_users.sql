-- +goose Up
-- +goose StatementBegin
INSERT INTO users (
    uuid, username, email, password
) VALUES (
    "fa455b13-0f39-4aff-ac4e-d24112ec76b6",
    "pm",
    "pm@pm.com",
    "$2a$14$7QHKumYyhchDIqOiAwolEeIxFBtrOa/efganjTlJYUf96zmIjvcJC"
);

INSERT INTO users (
    uuid, username, email, password
) VALUES (
    "c9568206-1ccf-4c9a-97eb-84bde7ad4605",
    "pa",
    "pa@pa.com",
    "$2a$14$G2C./eUg6L7MpTzp7o8/beLrloMlA2Jv/uuCf6e1U9anIXSvL9xsa"
);
-- +goose StatementEnd

-- +goose Down
DELETE FROM users;
