-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS student_subject (
    id BIGSERIAL,
	student_id  INT NOT NULL,
    subject_id   INT NOT NULL,
    marks  INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL,

	PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_subject;
-- +goose StatementEnd
