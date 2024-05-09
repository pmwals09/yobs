-- +goose Up
-- +goose StatementBegin
INSERT INTO opportunity_documents (
    opportunity_id,
    document_id
) VALUES (
    1, 1
);

INSERT INTO opportunity_documents (
    opportunity_id,
    document_id
) VALUES (
    2, 2
);

INSERT INTO opportunity_documents (
    opportunity_id,
    document_id
) VALUES (
    3, 3
);

INSERT INTO opportunity_documents (
    opportunity_id,
    document_id
) VALUES (
    4, 4
);

INSERT INTO opportunity_documents (
    opportunity_id,
    document_id
) VALUES (
    5, 1
);

INSERT INTO opportunity_documents (
    opportunity_id,
    document_id
) VALUES (
    6, 2
);
-- +goose StatementBegin

-- +goose Down
DELETE FROM opportunity_documents;
