-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS public.user (
    id serial PRIMARY KEY,
    login VARCHAR(128),
    password_hash VARCHAR(64)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE public.user;
-- +goose StatementEnd
