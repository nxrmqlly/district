-- +goose Up

ALTER TABLE users
DROP CONSTRAINT username_length;

ALTER TABLE users
ADD CONSTRAINT username_length CHECK (
    char_length(username) >= 3
    AND char_length(username) <= 24 -- reduced
);

-- +goose Down

ALTER TABLE users
DROP CONSTRAINT username_length;

ALTER TABLE users
ADD CONSTRAINT username_length CHECK (
    char_length(username) >= 3
    AND char_length(username) <= 40 -- old
);