-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE public.user ADD CONSTRAINT login_unique UNIQUE (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE public.user DROP CONSTRAINT login_unique;
-- +goose StatementEnd
