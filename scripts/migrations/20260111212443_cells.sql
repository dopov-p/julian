-- +goose Up
-- +goose StatementBegin
CREATE TABLE cells
(
    id                 text primary key,
    name               text      not null,
    kind               integer   not null default 0,
    content_type       integer   not null default 0,
    contents           jsonb,
    can_has_few_orders boolean   not null default false,
    created_at         timestamp not null,
    updated_at         timestamp not null,
    deleted_at         timestamp
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cells_name_unique
    ON cells(name)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cells;

DROP INDEX idx_cells_name_unique;
-- +goose StatementEnd
