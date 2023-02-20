-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL,
	class_id INT NOT NULL,
    first_name  TEXT NOT NULL,
    last_name TEXT NOT NULL,
    roll INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL,

	PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS students ;
-- +goose StatementEnd
