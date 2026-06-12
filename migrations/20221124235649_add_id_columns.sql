-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE public.balance ADD COLUMN id serial PRIMARY KEY;
ALTER TABLE public.withdrawal ADD COLUMN id serial PRIMARY KEY;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE public.balance DROP COLUMN id;
ALTER TABLE public.withdrawal DROP COLUMN id;
-- +goose StatementEnd
