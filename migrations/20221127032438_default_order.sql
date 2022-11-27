-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
alter table public.order alter column accrual set default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
