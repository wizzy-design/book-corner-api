-- +goose Up
INSERT INTO books (title, author, description, price, stock, image_url, status)
VALUES 
    ('Blackbeard', 'J.K. Rowling', 'A riveting story about pirates, treasure, and the high seas.', 4000, 5, 'https://images.unsplash.com/photo-1544947950-fa07a98d237f', 'published'),
    ('The Silent Forest', 'Elena Rostova', 'A gripping mystery set in the quiet pine wilderness of northern Europe.', 3500, 3, 'https://images.unsplash.com/photo-1512820790803-83ca734da794', 'published'),
    ('Drafts & Secrets', 'Morgan Lee', 'An unreleased thriller draft that is not yet ready for public readers.', 5000, 0, '', 'draft');

INSERT INTO tags (name, kind)
VALUES 
    ('Fantasy', 'genre'),
    ('Mystery', 'genre'),
    ('Thriller', 'genre');

INSERT INTO book_tags (book_id, tag_id)
VALUES 
    (
        (SELECT id FROM books WHERE title = 'Blackbeard' LIMIT 1),
        (SELECT id FROM tags WHERE name = 'Fantasy' LIMIT 1)
    ),
    (
        (SELECT id FROM books WHERE title = 'The Silent Forest' LIMIT 1),
        (SELECT id FROM tags WHERE name = 'Mystery' LIMIT 1)
    );

-- +goose Down
DELETE FROM books WHERE title IN ('Blackbeard', 'The Silent Forest', 'Drafts & Secrets');
DELETE FROM tags WHERE name IN ('Fantasy', 'Mystery', 'Thriller');
