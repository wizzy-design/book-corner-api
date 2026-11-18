-- +goose Up

CREATE TABLE quiz_questions(
    id SERIAL PRIMARY KEY,
    prompt TEXT NOT NULL, 
    position INT NOT NULL
);

CREATE TABLE quiz_options(
    id SERIAL PRIMARY KEY,
    question_id INT NOT NULL REFERENCES quiz_questions(id) ON DELETE CASCADE,
    label TEXT NOT NULL
);

CREATE TABLE option_tag_weights(
    option_id INT NOT NULL REFERENCES quiz_options(id) ON DELETE CASCADE,
    weight INT NOT NULL,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (option_id, tag_id)
);

-- +goose Down
DROP TABLE IF EXISTS option_tag_weights;
DROP TABLE IF EXISTS quiz_options;
DROP TABLE IF EXISTS quiz_questions;