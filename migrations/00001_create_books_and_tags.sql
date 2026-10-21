-- +goose Up
-- SQL in this section is executed when the migration is applied

CREATE TABLE books(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price INT NOT NULL,
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
    image_url TEXT NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    search tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(author, '')), 'B') ||
    setweight(to_tsvector('english', coalesce(description, '')), 'C')
    ) STORED
);

CREATE INDEX books_search_idx ON books USING GIN (search);

CREATE TABLE tags(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    kind VARCHAR(50) NOT NULL
);

CREATE TABLE book_tags(
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, tag_id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS book_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS books;