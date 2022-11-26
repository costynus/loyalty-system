-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE public.withdrawal ALTER COLUMN updated_at SET DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE public.withdrawal ALTER COLUMN updated_at DROP DEFAULT;
-- +goose StatementEnd
