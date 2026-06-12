-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE public.order ALTER COLUMN status SET DEFAULT 'NEW';
ALTER TABLE public.order ALTER COLUMN uploaded_at SET DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE public.order ALTER COLUMN status DROP DEFAULT;
ALTER TABLE public.order ALTER COLUMN uploaded_at DROP DEFAULT;
-- +goose StatementEnd
