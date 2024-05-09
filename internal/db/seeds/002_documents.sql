-- +goose Up
-- +goose StatementBegin
INSERT INTO documents (
    file_name,
    title,
    type,
    content_type,
    user_id
) VALUES (
    "Anonymized pmwals09 Resume.pdf",
    "My anonymous resume",
    "Resume",
    "application/pdf",
    1
);

INSERT INTO documents (
    file_name,
    title,
    type,
    content_type,
    user_id
) VALUES (
    "Anonymized pmwals09 Resume.pdf",
    "My anonymous resume",
    "Resume",
    "application/pdf",
    1
);

INSERT INTO documents (
    file_name,
    title,
    type,
    content_type,
    user_id
) VALUES (
    "Anonymized pmwals09 Resume.pdf",
    "new thing",
    "Communication",
    "application/pdf",
    1
);
-- +goose StatementEnd

-- +goose Down
DELETE FROM documents;
