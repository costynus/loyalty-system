-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE public.balance ADD CONSTRAINT user_id_unique UNIQUE (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE public.balance DROP CONSTRAINT user_id_unique;
-- +goose StatementEnd
